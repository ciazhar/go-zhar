package service

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"time"
)

type BasicService struct {
	repository repository.RedisRepository
}

func (b *BasicService) Get() (string, error) {
	return b.repository.Get()
}

func (b *BasicService) Set(value string, expiration time.Duration) error {
	return b.repository.Set(value, expiration)
}

func (b *BasicService) GetHash(field string) (string, error) {
	return b.repository.GetHash(field)
}

func (b *BasicService) SetHash(field string, value string) error {
	return b.repository.SetHash(field, value)
}

func (b *BasicService) SetHashTTL(field string, value string, ttl time.Duration) error {
	return b.repository.SetHashTTL(field, value, ttl)
}

func (b *BasicService) DeleteHash(field string) error {
	return b.repository.DeleteHash(field)
}

func NewBasicService(repository repository.RedisRepository) *BasicService {
	return &BasicService{
		repository: repository,
	}
}
