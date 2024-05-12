package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

func main() {
	socket, _, err := gows.NewClient(new(WebSocket), &gows.ClientOption{
		Addr: "ws://127.0.0.1:3000/connect",
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	if err != nil {
		log.Print(err.Error())
		return
	}
	go socket.ReadLoop()

	for {
		var text = ""
		fmt.Scanf("%s", &text)
		if strings.TrimSpace(text) == "" {
			continue
		}
		socket.WriteString(text)
	}
}

type WebSocket struct {
}

func (c *WebSocket) OnClose(socket *gows.Conn, err error) {
	fmt.Printf("onerror: err=%s\n", err.Error())
}

func (c *WebSocket) OnPong(socket *gows.Conn, payload []byte) {
}

func (c *WebSocket) OnOpen(socket *gows.Conn) {
	_ = socket.WriteString("hello, there is client")
}

func (c *WebSocket) OnPing(socket *gows.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *WebSocket) OnMessage(socket *gows.Conn, message *gows.Message) {
	defer message.Close()
	fmt.Printf("recv: %s\n", message.Data.String())
}
