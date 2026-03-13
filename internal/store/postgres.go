package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrShortURLExists   = errors.New("short url already exists")
	ErrShortURLNotFound = errors.New("short url not found")
)

type ShortURL struct {
	ID        string
	URL       string
	CreatedAt time.Time
}

type ShortURLStore struct {
	pool *pgxpool.Pool
}

func NewShortURLStore(ctx context.Context, databaseURL string) (*ShortURLStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &ShortURLStore{pool: pool}, nil
}

func (s *ShortURLStore) Close() {
	s.pool.Close()
}

func (s *ShortURLStore) Init(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS short_urls (
			id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return err
}

func (s *ShortURLStore) Create(ctx context.Context, id string, url string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO short_urls (id, url)
		VALUES ($1, $2)
	`, id, url)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrShortURLExists
	}

	return err
}

func (s *ShortURLStore) Get(ctx context.Context, id string) (ShortURL, error) {
	var shortURL ShortURL

	err := s.pool.QueryRow(ctx, `
		SELECT id, url, created_at
		FROM short_urls
		WHERE id = $1
	`, id).Scan(&shortURL.ID, &shortURL.URL, &shortURL.CreatedAt)
	if err == nil {
		return shortURL, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ShortURL{}, ErrShortURLNotFound
	}

	return ShortURL{}, err
}
