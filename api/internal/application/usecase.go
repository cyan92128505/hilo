package usecase

import (
	"context"
	"hilo-api/internal/domain/dto"
)

type IAuthCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
}

// IUserCase interface
type IUserCase interface {
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
}

// IMessageCase interface
type IMessageCase interface {
	SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*dto.MessageResponse, error)
	MarkAsRead(ctx context.Context, req *dto.MarkAsReadRequest) (*dto.MessageResponse, error)
	ListMessages(ctx context.Context, req *dto.ListMessagesRequest) (*dto.ListMessagesResponse, error)
	ListConversations(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListConversationsResponse, error)
}

// Set struct
type Set struct {
	AuthCase    IAuthCase
	UserCase    IUserCase
	MessageCase IMessageCase
}
