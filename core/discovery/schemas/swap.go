package schemas

import (
	"encoding/binary"
)

type SwapMessageType int // The message type for the swap message

const (
	Secret SwapMessageType = iota
)

type SwapMessage struct {
	Type    SwapMessageType `json:"type"`
	Payload []byte          `json:"payload"`
}

func (s *SwapMessage) Serialize() []byte {
	data := make([]byte, 32)

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(s.Type))
	copy(data[:2], buf)
	copy(data[2:], s.Payload)

	return data
}

func DeserializeSwapMessage(data []byte) *SwapMessage {
	return &SwapMessage{
		Type:    SwapMessageType(binary.BigEndian.Uint16(data[:2])),
		Payload: data[2:],
	}
}
