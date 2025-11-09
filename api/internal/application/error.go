package usecase

import (
	"errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrReceiverNotFound   = errors.New("receiver not found")
	ErrMessageNotFound    = errors.New("message not found")
)
