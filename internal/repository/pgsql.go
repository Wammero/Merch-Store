package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGRepo struct {
	mu   sync.Mutex
	pool *pgxpool.Pool
}

func New(connstr string) (*PGRepo, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	return &PGRepo{pool: pool}, nil
}

func (repo *PGRepo) Close() {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.pool != nil {
		repo.pool.Close()
	}
}
