package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"kafka-app/pkg/config"
)

type PostgresMessageRep struct {
	db *pgxpool.Pool
}

func NewMR(ctx context.Context, cfg *config.DB) (*PostgresMessageRep, error) {
	
	var err error
	connString := createConnString(cfg)

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
	err = pgInstance.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to connect db: %w", err)
	}
	
	if err = pgInstance.RunMigrations(ctx); err != nil {
		return nil, fmt.Errorf("fail to run migrations: %w", err)
	}

	return pgInstance, nil
}

func createConnString(cfg *config.DB) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}


func (pg *PostgresMessageRep) RunMigrations(ctx context.Context) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS message (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		read BOOLEAN NOT NULL DEFAULT false
	);
	`
	_, err := pg.db.Exec(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("unable to run migrations: %w", err)
	}
	return nil
}


func (pg *PostgresMessageRep) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE message SET read = true WHERE id = $1`
	_, err := pg.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to mark message as read: %w", err)
	}
	return nil
}