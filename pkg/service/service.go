package service

import (
	"io"

	"impulse-api/internal/entity"
	"impulse-api/pkg/repository"
)

type ZodiacApi interface {
	DataWorker(r io.Reader) (entity.Summary, error)
	GenerateToken(clientID int, clientSecret string) (string, error)
}

type Service struct {
	ZodiacApi
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ZodiacApi: NewWesternHoroscope(repo),
	}
}
