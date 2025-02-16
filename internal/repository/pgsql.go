package repository

import (
	"context"
	"errors"
	"sync"
	"time"

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
	config, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		return nil, err
	}

	// Настройки пула
	config.MaxConns = 50                       // Поддержка 50 соединений
	config.MinConns = 10                       // Минимум 10 соединений
	config.MaxConnLifetime = time.Minute * 5   // Соединение живёт 5 минут
	config.MaxConnIdleTime = time.Second * 30  // Закрывать соединения после 30 секунд простоя
	config.HealthCheckPeriod = time.Minute * 1 // Проверять соединения раз в минуту

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
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
