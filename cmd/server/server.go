package main

import (
	"fmt"
	"github.com/alzhyganshal/test-go-server-client/pkg/messages"
	"net"
	"sync"
)

const (
	ServerPort = "5555"
)

var freeTag int32
var conns map[byte]net.Conn
var rw sync.RWMutex

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		panic("Can't start server on port " + ServerPort + "! Error:" + err.Error())
	}
	defer l.Close()
	conns = make(map[byte]net.Conn)
	fmt.Println("Server started on " + ServerPort + " port...")
	// main loop waiting for client connections
	for {
		c, err := l.Accept()
		if err != nil {
			panic("Could not accept connection! " + err.Error())
		}
		go acceptConnection(c)
	}
}

/*
acceptConnection accepts new connection and manage actions
*/
func acceptConnection(c net.Conn) {
	// lock mutex before saving connection to our map
	rw.Lock()
	currentTag := byte(freeTag)
	conns[currentTag] = c
	ok := true
	// search for a free tag
	for ok {
		freeTag++
		_, ok = conns[byte(freeTag)]
	}
	rw.Unlock()
	defer func() {
		c.Close()
		rw.Lock()
		// delete connection tag from map
		delete(conns, currentTag)
		// assign free tag
		freeTag = int32(currentTag)
		rw.Unlock()
	}()
	b := make([]byte, 65536)
	for {
		// wait for incoming packet
		_, err := c.Read(b)
		if err != nil {
			fmt.Println("Could not read incoming message! " + err.Error())
			return
		}
		//fmt.Printf("Debug info: %X\n", b[:n])
		m := messages.UnpackMessage(b)
		// got init action: tag new client
		if m.Action == 0 {
			c.Write(messages.PackMessage(messages.Message{Action: 1, Address: currentTag}))
			fmt.Printf("New client with tag %d!\n", currentTag)
		}
		// got send action
		if m.Action == 2 {
			c.Write(messages.PackMessage(messages.Message{Action: 255}))
			// process broadcast message
			if m.Address == 255 {
				fmt.Printf("Broadcast message from %d: %s\n", m.From, m.Text)
				for _, v := range conns {
					// we should lock mutex when operate with other connections
					rw.Lock()
					v.Write(messages.PackMessage(m))
					rw.Unlock()
				}
			} else {
				// process personal message
				fmt.Printf("Incoming message from %d to %d: %s\n", m.From, m.Address, m.Text)
				recipient, ok := conns[m.Address]
				if !ok {
					fmt.Printf("Recipient %d is not available!", m.Address)
					continue
				}
				rw.Lock()
				recipient.Write(messages.PackMessage(m))
				rw.Unlock()
			}
		}
		// got ping-pong action
		if m.Action == 10 {
			c.Write(messages.PackMessage(messages.Message{Action: 11}))
		}
	}
}
