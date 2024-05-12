package main

import (
	"log"
	"net/http"

	gows "github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

func main() {
	upgrader := gows.NewUpgrader(&Handler{}, &gows.ServerOption{
		CheckUtf8Enabled: true,
		Recovery:         gows.Recovery,
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			return
		}
		go func() {
			socket.ReadLoop()
		}()
	})
	log.Panic(
		http.ListenAndServe(":8000", nil),
	)
}

type Handler struct {
	gows.BuiltinEventHandler
}

func (c *Handler) OnPing(socket *gows.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Handler) OnMessage(socket *gows.Conn, message *gows.Message) {
	defer message.Close()
	_ = socket.WriteMessage(message.Opcode, message.Bytes())
}
