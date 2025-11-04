package message

import "errors"

var (
	ErrReceiverNotFound = errors.New("receiver not found")
	ErrMessageNotFound  = errors.New("message not found")
)
