package usecase

import (
	"context"

	"kafka-app/internal/domain"
)

type MessageUsecase struct {
	repo domain.MessageRepository
}

func NewMessageUsecase(repo domain.MessageRepository) *MessageUsecase {
	return &MessageUsecase{repo: repo}
}

func (uc *MessageUsecase) Insert(ctx context.Context, text string) error {
	return uc.repo.Insert(ctx, text)
}

func (uc *MessageUsecase) GetUser(ctx context.Context, name string, limit int) ([]domain.Message, error) {
	return uc.repo.GetUser(ctx, name, limit)
}
