package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInsufficientFunds  = errors.New("insufficient funds for transfer")
	ErrReceiverNotFound   = errors.New("receiver not found")
	ErrTransactionFailed  = errors.New("transaction recording failed")
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

func (repo *PGRepo) GetPool() *pgxpool.Pool {
	// repo.mu.Lock()
	// defer repo.mu.Unlock()
	return repo.pool
}

func (repo *PGRepo) Close() {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.pool != nil {
		repo.pool.Close()
	}
}
