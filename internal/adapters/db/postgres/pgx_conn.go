package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"kafka-app/internal/domain"
	"kafka-app/pkg/config"
)

type PostgresMessageRep struct {
	db *pgxpool.Pool
}

func NewMR(ctx context.Context, cfg *config.DB) (*PostgresMessageRep, error) {
	var err error
  var connString string

	var pgInstance *PostgresMessageRep
	var pgOnce sync.Once

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = &PostgresMessageRep{db}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("fail to create connection: %w", err)
	}

	return pgInstance, nil
}

func (pg *PostgresMessageRep) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *PostgresMessageRep) Close() {
	pg.db.Close()
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
	query := `SELECT id, text FROM message WHERE name = $1 LIMIT $2`
	rows, err := pg.db.Query(ctx, query, name, limit)
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
