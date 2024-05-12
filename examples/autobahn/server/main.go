package main

import (
	"log"

	gows "github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

func main() {
	s1 := gows.NewServer(&Handler{Sync: true}, &gows.ServerOption{
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
		CheckUtf8Enabled: true,
		Recovery:         gows.Recovery,
	})

	s2 := gows.NewServer(&Handler{Sync: false}, &gows.ServerOption{
		ParallelEnabled: true,
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
		CheckUtf8Enabled: true,
		Recovery:         gows.Recovery,
	})

	s3 := gows.NewServer(&Handler{Sync: true}, &gows.ServerOption{
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: false,
			ClientContextTakeover: false,
		},
		CheckUtf8Enabled: true,
		Recovery:         gows.Recovery,
	})

	s4 := gows.NewServer(&Handler{Sync: false}, &gows.ServerOption{
		ParallelEnabled: true,
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: false,
			ClientContextTakeover: false,
		},
		CheckUtf8Enabled: true,
		Recovery:         gows.Recovery,
	})

	go func() {
		log.Panic(s1.Run(":8000"))
	}()

	go func() {
		log.Panic(s2.Run(":8001"))
	}()

	go func() {
		log.Panic(s3.Run(":8002"))
	}()

	log.Panic(s4.Run(":8003"))
}

type Handler struct {
	gows.BuiltinEventHandler
	Sync bool
}

func (c *Handler) OnPing(socket *gows.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Handler) OnMessage(socket *gows.Conn, message *gows.Message) {
	if c.Sync {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
		_ = message.Close()
	} else {
		socket.WriteAsync(message.Opcode, message.Bytes(), func(err error) { _ = message.Close() })
	}
}
