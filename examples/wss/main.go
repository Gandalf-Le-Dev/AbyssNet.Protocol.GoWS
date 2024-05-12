package main

import (
	"flag"
	"log"
	"path/filepath"

	gows "github.com/Gandalf-Le-Dev/abyssnet.protocol.gows"
)

var dir string

func init() {
	flag.StringVar(&dir, "d", "", "cert directory")
	flag.Parse()

	d, err := filepath.Abs(dir)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	dir = d
}

func main() {
	srv := gows.NewServer(new(Websocket), nil)

	// wss://www.gows.com:8443/
	if err := srv.RunTLS(":8443", dir+"/server.crt", dir+"/server.pem"); err != nil {
		log.Panicln(err.Error())
	}
}

type Websocket struct {
	gows.BuiltinEventHandler
}

func (c *Websocket) OnPing(socket *gows.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Websocket) OnMessage(socket *gows.Conn, message *gows.Message) {
	defer message.Close()
	_ = socket.WriteMessage(message.Opcode, message.Bytes())
}
