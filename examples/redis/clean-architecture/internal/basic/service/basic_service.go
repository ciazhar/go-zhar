package service

import (
	"github.com/ciazhar/go-zhar/examples/redis/clean-architecture/internal/basic/repository"
)

type BasicService interface {
	GetBasicHash(key string) (string, error)
	SetBasicHash(key string, value string) error
}

type basicService struct {
	repository repository.RedisRepository
}

func (b basicService) GetBasicHash(key string) (string, error) {
	return b.repository.GetBasicHash(key)
}

func (b basicService) SetBasicHash(key string, value string) error {
	return b.repository.SetBasicHash(key, value)
}

func NewBasicService(repository repository.RedisRepository) BasicService {
	return basicService{
		repository: repository,
	}
}
