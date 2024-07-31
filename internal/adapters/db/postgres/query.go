package postgres

import (
	"context"
	"fmt"

	"kafka-app/internal/domain"
)

func (pg *PostgresMessageRep) Close() {
	pg.db.Close()
}

func (pg *PostgresMessageRep) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *PostgresMessageRep) Insert(ctx context.Context, text string) error {
	query := `INSERT INTO message (text) VALUES ($1)`
	_, err := pg.db.Exec(ctx, query, text)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (pg *PostgresMessageRep) GetUser(ctx context.Context, name string, limit int) ([]domain.Message, error) {
	query := `SELECT id, text FROM message LIMIT $1`
	rows, err := pg.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fail to query users: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.Text); err != nil {
			return nil, fmt.Errorf("fail to scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return messages, nil
}
