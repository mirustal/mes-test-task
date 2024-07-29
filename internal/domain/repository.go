package domain

import (
	"context"
)

type MessageRepository interface {
	Insert(ctx context.Context, text string) error
	GetUser(ctx context.Context, name string, limit int) ([]Message, error)
}
