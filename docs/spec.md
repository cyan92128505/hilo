## 專案體驗流程

1. 註冊/登入
   POST /api/v1/auth/register
   POST /api/v1/auth/login
   → 回傳 JWT

2. 搜尋使用者
   GET /api/v1/users/search?q=username
   → 回傳使用者列表

3. 開啟聊天室
   POST /api/v1/rooms
   Body: {"participant_id": "user_b_id"}

4. WebSocket 連線
   GET /api/v1/ws?token=jwt
   → 建立 WebSocket 連線

5. 發送/接收訊息
   - 發送：透過 WebSocket 送 JSON
   - 接收：透過 WebSocket 收 JSON
   - 同時存入資料庫（持久化）