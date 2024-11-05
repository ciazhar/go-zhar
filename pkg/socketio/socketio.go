package socketio

import (
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"sync"
)

// ChannelHandler Define the function type that accepts *gosocketio.Channel as a parameter
type ChannelHandler func(c *gosocketio.Channel)

const SocketIORoomName = "digisar"

func Init(funcs ...ChannelHandler) *gosocketio.Server {
	var once sync.Once
	var server *gosocketio.Server

	once.Do(func() {
		server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	})

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")

		//join them to room
		c.Join(SocketIORoomName)

		for _, fn := range funcs {
			fn(c)
		}
	})

	return server
}
