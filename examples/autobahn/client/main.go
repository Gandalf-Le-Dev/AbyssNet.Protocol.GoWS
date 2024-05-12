package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

const remoteAddr = "127.0.0.1:9001"

func main() {
	const count = 517
	for i := 1; i <= count; i++ {
		testCase(true, i, "gows-client/sync")
	}
	for i := 1; i <= count; i++ {
		testCase(false, i, "gows-client/async")
	}
	updateReports()
}

func testCase(sync bool, id int, agent string) {
	var url = fmt.Sprintf("ws://%s/runCase?case=%d&agent=%s", remoteAddr, id, agent)
	var handler = &WebSocket{Sync: sync, onexit: make(chan struct{})}
	socket, _, err := gows.NewClient(handler, &gows.ClientOption{
		Addr:             url,
		CheckUtf8Enabled: true,
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	go socket.ReadLoop()
	<-handler.onexit
}

type WebSocket struct {
	Sync   bool
	onexit chan struct{}
}

func (c *WebSocket) OnOpen(socket *gows.Conn) {
	_ = socket.SetDeadline(time.Now().Add(30 * time.Second))
}

func (c *WebSocket) OnClose(socket *gows.Conn, err error) {
	c.onexit <- struct{}{}
}

func (c *WebSocket) OnPing(socket *gows.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *WebSocket) OnPong(socket *gows.Conn, payload []byte) {}

func (c *WebSocket) OnMessage(socket *gows.Conn, message *gows.Message) {
	if c.Sync {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
		_ = message.Close()
	} else {
		socket.WriteAsync(message.Opcode, message.Bytes(), func(err error) { _ = message.Close() })
	}
}

type updateReportsHandler struct {
	onexit chan struct{}
	gows.BuiltinEventHandler
}

func (c *updateReportsHandler) OnOpen(socket *gows.Conn) {
	_ = socket.SetDeadline(time.Now().Add(5 * time.Second))
}

func (c *updateReportsHandler) OnClose(socket *gows.Conn, err error) {
	c.onexit <- struct{}{}
}

func updateReports() {
	var url = fmt.Sprintf("ws://%s/updateReports?agent=gows/client", remoteAddr)
	var handler = &updateReportsHandler{onexit: make(chan struct{})}
	socket, _, err := gows.NewClient(handler, &gows.ClientOption{
		Addr:             url,
		CheckUtf8Enabled: true,
		PermessageDeflate: gows.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	go socket.ReadLoop()
	<-handler.onexit
}
