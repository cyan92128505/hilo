package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hilo-api/internal/domain/do"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *do.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at, read_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		msg.ID(),
		msg.SenderID(),
		msg.ReceiverID(),
		msg.Content(),
		msg.CreatedAt(),
		msg.ReadAt(),
	)
	return err
}

func (r *MessageRepository) FindByID(ctx context.Context, id uuid.UUID) (*do.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at, read_at
		FROM messages
		WHERE id = $1
	`

	var (
		msgID      uuid.UUID
		senderID   uuid.UUID
		receiverID uuid.UUID
		content    string
		createdAt  time.Time
		readAt     sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msgID, &senderID, &receiverID, &content, &createdAt, &readAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	var readAtPtr *time.Time
	if readAt.Valid {
		readAtPtr = &readAt.Time
	}

	return do.ReconstructMessage(msgID, senderID, receiverID, content, createdAt, readAtPtr), nil
}

func (r *MessageRepository) UpdateReadAt(ctx context.Context, id uuid.UUID, readAt time.Time) error {
	query := `
		UPDATE messages
		SET read_at = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, readAt, id)
	return err
}

func (r *MessageRepository) ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]*do.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at, read_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2)
		   OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, userA, userB, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*do.Message
	for rows.Next() {
		var (
			id         uuid.UUID
			senderID   uuid.UUID
			receiverID uuid.UUID
			content    string
			createdAt  time.Time
			readAt     sql.NullTime
		)

		if err := rows.Scan(&id, &senderID, &receiverID, &content, &createdAt, &readAt); err != nil {
			return nil, err
		}

		var readAtPtr *time.Time
		if readAt.Valid {
			readAtPtr = &readAt.Time
		}

		messages = append(messages, do.ReconstructMessage(id, senderID, receiverID, content, createdAt, readAtPtr))
	}

	return messages, rows.Err()
}

func (r *MessageRepository) ListUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*do.ConversationPreview, error) {
	// Get latest message from each conversation using ROW_NUMBER
	query := `
		WITH conversation_messages AS (
			-- Identify other user and rank messages by time
			SELECT 
				id, 
				sender_id, 
				receiver_id, 
				content, 
				created_at, 
				read_at,
				CASE 
					WHEN sender_id = $1 THEN receiver_id 
					ELSE sender_id 
				END as other_user_id,
				ROW_NUMBER() OVER (
					PARTITION BY CASE 
						WHEN sender_id = $1 THEN receiver_id 
						ELSE sender_id 
					END
					ORDER BY created_at DESC, id DESC
				) as rn
			FROM messages
			WHERE sender_id = $1 OR receiver_id = $1
		),
		latest_messages AS (
			-- Get only the latest message (rn = 1) from each conversation
			SELECT id, sender_id, receiver_id, content, created_at, read_at, other_user_id
			FROM conversation_messages
			WHERE rn = 1
		),
		unread_counts AS (
			-- Count unread messages from each user
			SELECT 
				sender_id as other_user_id,
				COUNT(*) as unread_count
			FROM messages
			WHERE receiver_id = $1 AND read_at IS NULL
			GROUP BY sender_id
		)
		SELECT 
			lm.id, lm.sender_id, lm.receiver_id, lm.content, lm.created_at, lm.read_at,
			u.id, u.email, u.password, u.username, u.created_at,
			COALESCE(uc.unread_count, 0) as unread_count
		FROM latest_messages lm
		JOIN users u ON u.id = lm.other_user_id
		LEFT JOIN unread_counts uc ON uc.other_user_id = lm.other_user_id
		ORDER BY lm.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previews []*do.ConversationPreview
	for rows.Next() {
		var (
			msgID         uuid.UUID
			senderID      uuid.UUID
			receiverID    uuid.UUID
			content       string
			msgCreatedAt  time.Time
			readAt        sql.NullTime
			userID        uuid.UUID
			email         string
			password      string
			username      string
			userCreatedAt time.Time
			unreadCount   int
		)

		if err := rows.Scan(
			&msgID, &senderID, &receiverID, &content, &msgCreatedAt, &readAt,
			&userID, &email, &password, &username, &userCreatedAt,
			&unreadCount,
		); err != nil {
			return nil, err
		}

		var readAtPtr *time.Time
		if readAt.Valid {
			readAtPtr = &readAt.Time
		}

		previews = append(previews, &do.ConversationPreview{
			OtherUser:   do.ReconstructUser(userID, email, password, username, userCreatedAt),
			LastMessage: do.ReconstructMessage(msgID, senderID, receiverID, content, msgCreatedAt, readAtPtr),
			UnreadCount: unreadCount,
		})
	}

	return previews, rows.Err()
}
