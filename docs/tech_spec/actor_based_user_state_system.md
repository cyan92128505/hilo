# Actor-based user state system
README + architecture + full Go example (single-file demo)
## 使用說明：這是一個示範性單檔實作（可拆分為多檔）。
## 需要的套件：
- go get github.com/go-redis/redis/v9
- go get github.com/gorilla/websocket
- go get github.com/jackc/pgx/v5
在真實專案中請把各部份拆成多個檔案與 package，增加錯誤處理與監控。


## 架構圖（Mermaid）

```mermaid
flowchart LR
  subgraph AppServer
    WS[WebSocket]
    ActorMgr[Actor Manager]
    DB[(PostgreSQL)]
    Redis[(Redis)]
  end

  Client -->|WS events| WS
  WS -->|On message| ActorMgr
  ActorMgr -->|Update: messages table| DB
  ActorMgr -->|Update: conversation summary| Redis
  DB -->|Optional read fallback| ActorMgr
  Redis -->|pub/sub cross-node| ActorMgr
```

### 說明：
- Client (App) 與 AppServer 維持 WebSocket 連線。
- 每個 user_id 有一個 Actor（goroutine + channel），由 Actor Manager 建立/回收。
- 發訊息時：寫入 PostgreSQL（durable），然後將相關事件送給 receiver 的 Actor，以更新 Redis 的 conversation summary（並推送即時通知）。
- 查詢會話列表直接從 Redis 讀取（ZSET + HASH），若快取失效則 fallback 到 DB 並回填快取。
- Redis pub/sub 用於多 AppServer 節點間的即時狀態同步（例如多台應用伺服器）。


```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

// 簡化的設定
var (
	redisAddr = "localhost:6379"
	pgConnStr = "postgres://user:pass@localhost:5432/dbname"
)

// 事件定義
type EventType string

const (
	EventMessageReceived EventType = "MessageReceived"
	EventMessageRead     EventType = "MessageRead"
	EventUserOnline      EventType = "UserOnline"
	EventUserOffline     EventType = "UserOffline"
)

type UserEvent struct {
	Type EventType
	Data any
}

// ActorState 儲存每個使用者的必要狀態
type ActorState struct {
	UserID    string
	Online    bool
	LastSeen  time.Time
	Unread    map[string]int64 // other_user_id -> count
	mutex     sync.Mutex       // 內部狀態保護（actor 內通常不需）
}

// UserActor 是每個使用者的 actor
type UserActor struct {
	UserID string
	Events chan UserEvent
	Stop   chan struct{}

	redis *redis.Client
	pg    *pgx.Conn

	state ActorState
}

func NewUserActor(userID string, r *redis.Client, pg *pgx.Conn) *UserActor {
	return &UserActor{
		UserID: userID,
		Events: make(chan UserEvent, 256),
		Stop:   make(chan struct{}),
		redis:  r,
		pg:     pg,
		state: ActorState{
			UserID:   userID,
			Online:   false,
			LastSeen: time.Now(),
			Unread:   make(map[string]int64),
		},
	}
}

func (a *UserActor) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("actor %s started", a.UserID)
	idleTimer := time.NewTimer(0)
	if !idleTimer.Stop() {
		<-idleTimer.C
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("actor %s context done", a.UserID)
			return
		case <-a.Stop:
			log.Printf("actor %s stop", a.UserID)
			return
		case ev := <-a.Events:
			a.handleEvent(ev)
		case <-idleTimer.C:
			// 可做定期 flush 或自我回收
			// 先不回收以示範
		}
	}
}

func (a *UserActor) handleEvent(ev UserEvent) {
	switch ev.Type {
	case EventMessageReceived:
		m := ev.Data.(MessagePayload)
		// 1) 將訊息寫到 DB (簡化示範, 實務要非同步或用 queue)
		if err := insertMessage(a.pg, m); err != nil {
			log.Printf("db insert error: %v", err)
		}
		// 2) 更新 Redis summary for receiver
		a.updateRedisOnMessage(m)
	case EventMessageRead:
		r := ev.Data.(ReadPayload)
		a.clearUnreadInRedis(r.OtherUserID)
	case EventUserOnline:
		a.setOnline(true)
	case EventUserOffline:
		a.setOnline(false)
	default:
		log.Printf("unknown event: %v", ev.Type)
	}
}

// MessagePayload 與 ReadPayload 是事件攜帶的資料
type MessagePayload struct {
	MessageID   string
	SenderID    string
	ReceiverID  string
	Content     string
	CreatedAt   time.Time
}

type ReadPayload struct {
	OtherUserID string
}

// insertMessage 示例
func insertMessage(conn *pgx.Conn, m MessagePayload) error {
	// 這只是同步示範；實務上可能用 background worker 或 queue
	_, err := conn.Exec(context.Background(),
		`INSERT INTO messages (id, sender_id, receiver_id, content, created_at) VALUES ($1,$2,$3,$4,$5)`,
		m.MessageID, m.SenderID, m.ReceiverID, m.Content, m.CreatedAt)
	return err
}

// Redis key convention
// conversation:index:<user_id> -> ZSET sorted by last_message_at, member = other_user_id
// conversation:summary:<user_id>:<other_user_id> -> HASH { last_message, last_message_at, unread }

func conversationIndexKey(userID string) string { return fmt.Sprintf("conversation:index:%s", userID) }
func conversationSummaryKey(userID, other string) string {
	return fmt.Sprintf("conversation:summary:%s:%s", userID, other)
}

func (a *UserActor) updateRedisOnMessage(m MessagePayload) {
	// 更新 receiver 的 conversation summary
	other := m.SenderID
	key := conversationSummaryKey(m.ReceiverID, other)
	zkey := conversationIndexKey(m.ReceiverID)

	// HMSET summary
	t := m.CreatedAt.Unix()
	// 使用 pipeline
	pipe := a.redis.Pipeline()
	pipe.HSet(context.Background(), key, map[string]any{
		"last_message":    m.Content,
		"last_message_at": t,
	})
	pipe.HIncrBy(context.Background(), key, "unread", 1)
	pipe.ZAdd(context.Background(), zkey, &redis.Z{Score: float64(t), Member: other})
	if _, err := pipe.Exec(context.Background()); err != nil {
		log.Printf("redis pipeline error: %v", err)
	}
}

func (a *UserActor) clearUnreadInRedis(other string) {
	key := conversationSummaryKey(a.UserID, other)
	pipe := a.redis.Pipeline()
	pipe.HSet(context.Background(), key, "unread", 0)
	if _, err := pipe.Exec(context.Background()); err != nil {
		log.Printf("redis clear unread error: %v", err)
	}
}

func (a *UserActor) setOnline(online bool) {
	a.state.mutex.Lock()
	a.state.Online = online
	a.state.LastSeen = time.Now()
	a.state.mutex.Unlock()
	// 更新 Redis 的線上集合
	if online {
		a.redis.SAdd(context.Background(), "online_users", a.UserID)
	} else {
		a.redis.SRem(context.Background(), "online_users", a.UserID)
	}
}

// ActorManager 管理所有 user actors
type ActorManager struct {
	actors sync.Map // map[userID]*UserActor
	redis  *redis.Client
	pg     *pgx.Conn
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewActorManager(r *redis.Client, pg *pgx.Conn) *ActorManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ActorManager{redis: r, pg: pg, ctx: ctx, cancel: cancel}
}

func (am *ActorManager) GetOrCreate(userID string) *UserActor {
	if v, ok := am.actors.Load(userID); ok {
		return v.(*UserActor)
	}
	actor := NewUserActor(userID, am.redis, am.pg)
	am.actors.Store(userID, actor)
	am.wg.Add(1)
	go actor.Run(am.ctx, &am.wg)
	return actor
}

func (am *ActorManager) StopAll() {
	am.cancel()
	am.wg.Wait()
}

// WebSocket handler 範例
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(am *ActorManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 簡化：userID 由 query 取得，實務請用 auth token
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user_id", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade err", err)
			return
		}
		defer conn.Close()

		actor := am.GetOrCreate(userID)
		// 標示在線上
		actor.Events <- UserEvent{Type: EventUserOnline}

		// 簡化的讀 loop
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("ws read err", err)
				actor.Events <- UserEvent{Type: EventUserOffline}
				return
			}
			// 這裡假設收到 JSON {type: "send", to: "other", content: "..."}
			// 為簡潔起見直接模擬成固定 payload
			payload := MessagePayload{
				MessageID:  fmt.Sprintf("msg-%d", time.Now().UnixNano()),
				SenderID:   userID,
				ReceiverID: "other-user-id-1",
				Content:    string(msg),
				CreatedAt:  time.Now(),
			}
			// 寫事件給 receiver 的 actor
			recActor := am.GetOrCreate(payload.ReceiverID)
			recActor.Events <- UserEvent{Type: EventMessageReceived, Data: payload}
		}
	}
}

func main() {
	// 連 Redis
	r := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()
	if err := r.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	// 連 Postgres（示範用單一連線）	
	pg, err := pgx.Connect(ctx, pgConnStr)
	if err != nil {
		log.Fatalf("pg connect: %v", err)
	}
	defer pg.Close(ctx)

	am := NewActorManager(r, pg)
	defer am.StopAll()

	http.HandleFunc("/ws", wsHandler(am))
	log.Println("server start :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```


## 注意與延伸：
1) Production 建議：actor pool 與 eviction 策略（閒置回收），避免大量 idle goroutine。
2) 多節點：使用 Redis pub/sub 或 Redis Streams 在不同 AppServer 間同步狀態。當收到 pub/sub 訊息時，對應 actor 可收到事件（或若 actor 不存在，將快取更新）。
3) durability：訊息寫入 DB 才是最終資料來源。Redis 僅作快取，故需設計 fallback 與回填策略。
4) 錯誤處理：示範程式省略大量錯誤處理、重試、監控、metrics、死信處理等，真實系統需要完善設計。
5) 性能優化：可用批次寫 DB、使用 background worker 處理 heavy I/O，actor 只下發事件到 worker queue。
