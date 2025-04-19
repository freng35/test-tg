package store

import (
	"context"
	"database/sql"
)

type Repository interface {
	AddUser(ctx context.Context, telegramID, channelID int64) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AddUser(ctx context.Context, telegramID, channelID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users(telegram_id, channel_id) 
		VALUES($1, $2) 
		ON CONFLICT (telegram_id) DO NOTHING`,
		telegramID,
		channelID,
	)
	return err
}
