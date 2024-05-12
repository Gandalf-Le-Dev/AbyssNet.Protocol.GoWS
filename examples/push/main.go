package main

import (
	"bufio"
	"log"
	"net"
	"net/http"

	gows "github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

func main() {
	var app = gows.NewServer(new(Handler), nil)

	app.OnRequest = func(conn net.Conn, br *bufio.Reader, r *http.Request) {
		socket, err := app.GetUpgrader().UpgradeFromConn(conn, br, r)
		if err != nil {
			log.Print(err.Error())
			return
		}
		var channel = make(chan []byte, 8)
		var closer = make(chan struct{})
		socket.Session().Store("channel", channel)
		socket.Session().Store("closer", closer)
		go socket.ReadLoop()
		go func() {
			for {
				select {
				case p := <-channel:
					_ = socket.WriteMessage(gows.OpcodeText, p)
				case <-closer:
					return
				}
			}
		}()
	}

	log.Fatalf("%v", app.Run(":8000"))
}

type Handler struct {
	gows.BuiltinEventHandler
}

func (c *Handler) getSession(socket *gows.Conn, key string) any {
	v, _ := socket.Session().Load(key)
	return v
}

func (c *Handler) Send(socket *gows.Conn, payload []byte) {
	var channel = c.getSession(socket, "channel").(chan []byte)
	select {
	case channel <- payload:
	default:
		return
	}
}

func (c *Handler) OnClose(socket *gows.Conn, err error) {
	var closer = c.getSession(socket, "closer").(chan struct{})
	closer <- struct{}{}
}

func (c *Handler) OnMessage(socket *gows.Conn, message *gows.Message) {
	defer message.Close()
	_ = socket.WriteMessage(message.Opcode, message.Bytes())
}
