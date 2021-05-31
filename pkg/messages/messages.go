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
	s := ""
	for i := 3; i < int(length)+3; i++ {
		s += string((*message)[i])
	}
	m.Text = s
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
