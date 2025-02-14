package service

import "Merch-Store/internal/repository"

type Service struct {
	repo *repository.PGRepo
}

func New(repo *repository.PGRepo) *Service {
	return &Service{repo: repo}
}
