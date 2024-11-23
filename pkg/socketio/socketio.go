package socketio

import (
	"context"
	"github.com/ciazhar/go-start-small/pkg/logger"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"sync"
)

// ChannelHandler Define the function type that accepts *gosocketio.Channel as a parameter
type ChannelHandler func(c *gosocketio.Channel)

func Init(roomName string, funcs ...ChannelHandler) *gosocketio.Server {
	var once sync.Once
	var server *gosocketio.Server

	once.Do(func() {
		server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	})

	err := server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")

		//join them to room
		err := c.Join(roomName)
		if err != nil {
			logger.LogFatal(context.Background(), err, "failed to join room", nil)
		}

		for _, fn := range funcs {
			fn(c)
		}
	})
	if err != nil {
		logger.LogFatal(context.Background(), err, "failed to initialize socketio", nil)
	}

	return server
}
