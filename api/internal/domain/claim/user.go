package claim

import (
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
)

var (
	// Now variable
	Now = time.Now
)

// ClaimsOption interface
type ClaimsOption interface {
	Apply(*User)
}

// WithUserID method
func WithUserID(id string) ClaimsOption {
	return withUserID{id: id}
}

type withUserID struct {
	id string
}

// Apply method
func (w withUserID) Apply(c *User) {
	c.UserID = w.id
}

// WithPermissions method
func WithPermissions(permission ...string) ClaimsOption {
	return withPermissions{permissions: permission}
}

type withPermissions struct {
	permissions []string
}

// Apply method
func (w withPermissions) Apply(c *User) {
	c.Permissions = w.permissions
}

// NewUser method
func NewUser(claims *jwt.Claims, options ...ClaimsOption) *User {
	user := &User{
		Claims: claims,
	}

	for _, option := range options {
		option.Apply(user)
	}

	return user
}

// User type
type User struct {
	UserID      string   `json:"user_id,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	*jwt.Claims
}

func (c *User) ExpiresAfter(d time.Duration) {
	c.Expiry = jwt.NewNumericDate(Now().Add(d))
}

func (c *User) GetExpiresAfter() *jwt.NumericDate {
	return c.Expiry
}
