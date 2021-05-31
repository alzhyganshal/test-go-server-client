package main

import (
	"fmt"
	"github.com/alzhyganshal/test-go-server-client/pkg/messages"
	"net"
)

const (
	SERVER_PORT = "5555"
)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", SERVER_PORT))
	if err != nil {
		panic("Can't start server on port " + SERVER_PORT + "! Error:" + err.Error())
	}
	defer l.Close()
	fmt.Println("Server started on " + SERVER_PORT + " port...")
	for {
		c, err := l.Accept()
		if err != nil {
			panic("Could not accept connection! " + err.Error())
		}
		go acceptConnection(c)
	}
}

func acceptConnection(c net.Conn) {
	b := make([]byte, 65536)
	n, err := c.Read(b)
	if err != nil {
		fmt.Println("Could not read incoming message! " + err.Error())
	}
	m := messages.UnpackMessage(&b)
	fmt.Printf("Incoming message to %d: %s\n", m.Address, m.Text)
	fmt.Printf("%X", b[:n])
	c.Write(messages.PackMessage(messages.Message{
		Text:    "Success",
		Address: 0,
	}))
}
