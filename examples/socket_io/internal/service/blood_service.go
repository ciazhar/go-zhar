package service

import (
    "context"
    "github.com/ciazhar/go-start-small/examples/socket_io/internal/model"
    gosocketio "github.com/graarh/golang-socketio"
)

type BloodService struct {
    latestBloodQueue chan *model.BloodAvailability
}

func NewBloodService(queueSize int) *BloodService {
    return &BloodService{
        latestBloodQueue: make(chan *model.BloodAvailability, queueSize),
    }
}

func (s *BloodService) BroadcastLatestBlood(ctx context.Context, blood *model.BloodAvailability) error {
    select {
    case s.latestBloodQueue <- blood:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (s *BloodService) ListenLatestBlood(server *gosocketio.Channel) {
    for blood := range s.latestBloodQueue {
        server.BroadcastTo("digisar", blood.ApplicationID+".blood.latest", blood)
    }
}