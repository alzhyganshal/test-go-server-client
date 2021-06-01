package main

import (
	"bufio"
	"fmt"
	"github.com/alzhyganshal/test-go-server-client/pkg/messages"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ServerPort = "5555"
	ServerName = "localhost"
)

var myTag byte
var rw sync.RWMutex

/*
sendAndReceive sends message to a server and waits for answer
*/
func sendAndReceive(c net.Conn, m messages.Message) (messages.Message, error) {
	// lock mutex to avoid concurrent connections
	rw.Lock()
	// unlock mutex in the end
	defer rw.Unlock()
	// pack message in bytes array and send it
	if _, err := c.Write(messages.PackMessage(m)); err != nil {
		return messages.Message{}, err
	}
	a := make([]byte, 65536)
	// wait for server reply and parse it
	if _, err := bufio.NewReader(c).Read(a); err != nil {
		return messages.Message{}, err
	}
	return messages.UnpackMessage(a), nil
}

/*
getCommands waits for user input and sends message or quit accordingly to command
*/
func getCommands(c net.Conn) {
	for {
		// reading from console
		r := bufio.NewReader(os.Stdin)
		t, err := r.ReadString('\n')
		if err != nil {
			panic("Could not read input! " + err.Error())
		}
		// exit if "quit" command was entered
		if strings.Contains(t, "quit") {
			os.Exit(0)
		}
		// send message action
		if strings.Contains(t, "send") {
			var addr int64
			// split command by space symbol
			sendString := strings.Split(t, " ")
			// asterisk means broadcast message
			if sendString[1] == "*" {
				addr = 255
			} else {
				addr, err = strconv.ParseInt(sendString[1], 10, 8)
				if err != nil {
					fmt.Println("Could not parse address tag: " + err.Error())
					continue
				}
			}
			// form message body
			message := messages.Message{
				Action:  2,
				From:    myTag,
				Text:    strings.Split(t, "\"")[1],
				Address: byte(addr),
			}
			// send message and receive acknowledge
			m, err := sendAndReceive(c, message)
			if err != nil {
				panic(err)
			}
			if m.Action == 255 {
				fmt.Println("Successfully sent!")
				continue
			}
			fmt.Println("Server did not accept message!")
			continue
		}
	}
}

func main() {
	c, err := net.Dial("tcp", ServerName+":"+ServerPort)
	defer c.Close()
	if err != nil {
		panic("Could not connect to server " + ServerName + ":" + ServerPort + "! " + err.Error())
	}
	// First we need an assigned id by server. Sending init action.
	m, err := sendAndReceive(c, messages.Message{Action: 0})
	if err != nil {
		panic(err)
	}
	myTag = m.Address
	fmt.Printf("Successfully connected to server %s:%s. Your id is %d.\n", ServerName, ServerPort, myTag)
	fmt.Println("Input message in following format: send <address> \"<message>\". Set address as * to broadcast message.")
	fmt.Println("Exit client by entering following command: quit")
	// run input console waiting as go routine
	go getCommands(c)

	// send and wait ping-pong messages to server
	for {
		m, err := sendAndReceive(c, messages.Message{Action: 10})
		if err != nil {
			panic(err)
		}
		// got message from other client
		if m.Action == 2 {
			fmt.Printf("Message from %d: %s\n", m.From, m.Text)
		}
		// some waiting of 0.5 second
		time.Sleep(500 * time.Millisecond)
	}
}
