package main

import (
	"bufio"
	"fmt"
	"github.com/alzhyganshal/test-go-server-client/pkg/messages"
	"net"
	"os"
	"strings"
)

const (
	SERVER_PORT = "5555"
	SERVER_NAME = "localhost"
)

func main() {
	fmt.Println("Input message in following format: send <address> <message>. Set address as * to broadcast message.")
	for {
		c, err := net.Dial("tcp", SERVER_NAME+":"+SERVER_PORT)
		if err != nil {
			panic("Could not connect to server " + SERVER_NAME + ":" + SERVER_PORT + "! " + err.Error())
		}
		defer c.Close()
		r := bufio.NewReader(os.Stdin)
		t, err := r.ReadString('\n')
		if err != nil {
			panic("Could not read input! " + err.Error())
		}
		if strings.Contains(t, "quit") {
			os.Exit(0)
		}
		if strings.Contains(t, "send") {
			message := messages.Message{
				Text:    strings.Split(t, " ")[2],
				Address: 0,
			}
			_, err = c.Write(messages.PackMessage(message))
			if err != nil {
				panic("Could not send message to server! " + err.Error())
			}
			a := make([]byte, 2048)
			_, err = bufio.NewReader(c).Read(a)
			aMessage := messages.UnpackMessage(&a)
			fmt.Println(aMessage.Text)
		}
	}
}
