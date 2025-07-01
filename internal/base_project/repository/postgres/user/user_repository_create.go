package user

import (
	"context"
	"log"
)

func (u userRepository) CreateUser(ctx context.Context, requestID string) error {
	log.Printf("[request_id=%s] inserting user to DB", requestID)
	// Simulasi insert ke database
	return nil
}
