package service

import (
	"io"

	"impulse-api/internal/entity"
	"impulse-api/pkg/repository"
)

type ZodiacApi interface {
	DataWorkerWithoutTime(r io.Reader, sex string) (entity.Summary, error)
	DataWorkerWithTime(r io.Reader) (entity.Summary, error)
	GenerateToken(clientID int) (string, error)
}

type Service struct {
	ZodiacApi
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ZodiacApi: NewWesternHoroscope(repo),
	}
}
