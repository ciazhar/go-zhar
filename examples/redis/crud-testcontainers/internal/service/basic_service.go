package service

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"time"
)

type BasicService interface {
	Get() (string, error)
	Set(value string, expiration time.Duration) error
	GetHash(field string) (string, error)
	SetHash(field string, value string) error
	SetHashTTL(field string, value string, ttl time.Duration) error
	DeleteHash(field string) error
}

type basicService struct {
	repository repository.RedisRepository
}

func (b basicService) Get() (string, error) {
	return b.repository.Get()
}

func (b basicService) Set(value string, expiration time.Duration) error {
	return b.repository.Set(value, expiration)
}

func (b basicService) GetHash(field string) (string, error) {
	return b.repository.GetHash(field)
}

func (b basicService) SetHash(field string, value string) error {
	return b.repository.SetHash(field, value)
}

func (b basicService) SetHashTTL(field string, value string, ttl time.Duration) error {
	return b.repository.SetHashTTL(field, value, ttl)
}

func (b basicService) DeleteHash(field string) error {
	return b.repository.DeleteHash(field)
}

func NewBasicService(repository repository.RedisRepository) BasicService {
	return basicService{
		repository: repository,
	}
}
