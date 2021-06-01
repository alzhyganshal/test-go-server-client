// Internal package for client-server app test task

package messages

import (
	"encoding/binary"
)

type Message struct {
	Text    string
	Address byte
	// Action thing can be on of the following values:
	// 0 - init action - client sends this action first to get assigned tag from server
	// 1 - tag action - server sends new tag after receiving init action from client
	// 2 - send action - messaging action
	// 10 - ping action - client sends this keep alive packets
	// 11 - pong action - server responds to ping action from client
	Action byte
	From   byte
}

/*
UnpackMessage parses bytes inside packet to a Message struct
*/
func UnpackMessage(message []byte) Message {
	var m Message
	// extract message length from first two bytes in Big Endian order
	length := binary.BigEndian.Uint16(message[:2])
	m.Action = message[2]
	// unpack send action
	if m.Action == 2 {
		m.From = message[3]
		m.Address = message[4]
		m.Text = string(message[5 : length+2])
		return m
	}
	// tag action got address field
	if m.Action == 1 {
		m.Address = message[3]
		return m
	}
	// all other actions does not have any additional content
	return m
}

/*
PackMessage forms bytes array from Message struct
*/
func PackMessage(m Message) []byte {
	// initialize 2 byte array for length
	b := make([]byte, 2)
	// tag action
	if m.Action == 1 {
		// set length to two bytes
		b[1] = 2
		b = append(b, m.Action)
		b = append(b, m.Address)
		return b
	}
	// send action
	if m.Action == 2 {
		// pack text length + 3 bytes (action, from, to) into bytes
		binary.BigEndian.PutUint16(b, uint16(len(m.Text)+3))
		tb := []byte(m.Text)
		b = append(b, m.Action)
		b = append(b, m.From)
		b = append(b, m.Address)
		b = append(b, tb...)
		return b
	}
	// all other actions does not have any content
	b[1] = 1
	b = append(b, m.Action)
	return b
}
