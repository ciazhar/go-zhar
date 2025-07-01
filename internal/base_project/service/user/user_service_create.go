package user

import (
	"context"
	"log"
)

func (u userService) CreateUser(ctx context.Context, requestID string) error {
	log.Printf("[request_id=%s] processing CreateUser in service", requestID)
	return u.repo.CreateUser(ctx, requestID)
}
