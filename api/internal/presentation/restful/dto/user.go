package dto

import (
	"hilo-api/internal/domain/do"
	"time"
)

// UserResponse represents a user
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// FromDomain converts domain user to DTO
func (u *UserResponse) FromDomain(user *do.User) {
	u.ID = user.ID().String()
	u.Email = user.Email()
	u.Username = user.Username()
	u.CreatedAt = user.CreatedAt()
}

// ListUsersRequest represents list users request
type ListUsersRequest struct {
	Limit  int `form:"limit" binding:"required,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// ListUsersResponse represents list users response
type ListUsersResponse struct {
	Users []*UserResponse `json:"users"`
	Total int             `json:"total"`
}
