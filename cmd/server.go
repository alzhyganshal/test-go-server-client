package main

import (
	"fmt"
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

}
