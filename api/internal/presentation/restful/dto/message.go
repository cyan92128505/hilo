package dto

import (
	"hilo-api/internal/domain/do"
	"time"
)

// SendMessageRequest represents send message request
type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" binding:"required,uuid"`
	Content    string `json:"content" binding:"required,max=5000"`
}

// MarkAsReadRequest represents mark as read request
type MarkAsReadRequest struct {
	MessageID string `json:"message_id" binding:"required,uuid"`
}

// MessageResponse represents a single message
type MessageResponse struct {
	ID         string     `json:"id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
}

// FromDomain converts domain message to DTO
func (m *MessageResponse) FromDomain(msg *do.Message) {
	m.ID = msg.ID().String()
	m.SenderID = msg.SenderID().String()
	m.ReceiverID = msg.ReceiverID().String()
	m.Content = msg.Content()
	m.CreatedAt = msg.CreatedAt()
	m.ReadAt = msg.ReadAt()
}

// ListMessagesRequest represents list messages request
type ListMessagesRequest struct {
	OtherUserID string `form:"user_id" binding:"required,uuid"`
	Limit       int    `form:"limit" binding:"required,min=1,max=100"`
	Offset      int    `form:"offset" binding:"omitempty,min=0"`
}

// ListMessagesResponse represents list messages response
type ListMessagesResponse struct {
	Messages []*MessageResponse `json:"messages"`
	Total    int                `json:"total"`
}

// ConversationPreviewResponse represents a conversation preview
type ConversationPreviewResponse struct {
	OtherUser   *UserResponse    `json:"other_user"`
	LastMessage *MessageResponse `json:"last_message"`
	UnreadCount int              `json:"unread_count"`
}

// FromDomain converts domain conversation preview to DTO
func (c *ConversationPreviewResponse) FromDomain(preview *do.ConversationPreview) {
	c.OtherUser = &UserResponse{}
	c.OtherUser.FromDomain(preview.OtherUser)

	c.LastMessage = &MessageResponse{}
	c.LastMessage.FromDomain(preview.LastMessage)

	c.UnreadCount = preview.UnreadCount
}

// ListConversationsResponse represents list conversations response
type ListConversationsResponse struct {
	Conversations []*ConversationPreviewResponse `json:"conversations"`
	Total         int                            `json:"total"`
}
