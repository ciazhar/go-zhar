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

func (b *BasicService) Delete() error {
	return b.repository.Delete()
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

func (b *BasicService) GetList() ([]string, error) {
	return b.repository.GetList()
}

func (b *BasicService) SetList(list []string) error {
	return b.repository.SetList(list)
}

func (b *BasicService) DeleteList(value string) error {
	return b.repository.DeleteList(value)
}

func NewBasicService(repository repository.RedisRepository) *BasicService {
	return &BasicService{
		repository: repository,
	}
}
