package core

import (
	"fmt"
	"slices"
)

const (
	MaxMessageSize      = 127
	ClientMask     byte = 0b00000000
	ServerMask     byte = 0b10000000

	MessageTypeReset      byte = 0
	MessageTypeIncrement  byte = 1
	MessageTypeGetCurrent byte = 2
	MessageTypeGetBest    byte = 3
)

type Transportable interface {
	Format() (Message, error)
}

type Message struct {
	Tag   byte
	Value []byte
}

func (m Message) TLV() []byte {
	l := len(m.Value)
	if l > MaxMessageSize {
		panic(fmt.Sprintf("Value bigger than max size: %d", l))
	}
	tlv := make([]byte, l+2)
	tlv = append(tlv, m.Tag)
	tlv = append(tlv, byte(l))
	return slices.Concat(tlv, m.Value)
}
