package messages

import (
	"encoding/binary"
)

type Message struct {
	Text    string
	Address byte
}

func UnpackMessage(message *[]byte) Message {
	var m Message
	length := binary.BigEndian.Uint16((*message)[:2])
	m.Address = (*message)[2]
	m.Text = string((*message)[3 : length+2])
	return m
}

func PackMessage(m Message) []byte {
	lb := make([]byte, 2)
	binary.BigEndian.PutUint16(lb, uint16(len(m.Text)))
	tb := []byte(m.Text)
	lb = append(lb, m.Address)
	lb = append(lb, tb...)
	return lb
}
