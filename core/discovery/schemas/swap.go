package schemas

import (
	"encoding/base64"
	"encoding/binary"
)

type SwapMessageType int // The message type for the swap message

const (
	Secret SwapMessageType = iota
	ContractCreated
)

type SwapMessage struct {
	Type    SwapMessageType `json:"type"`
	Payload []byte          `json:"payload"`
}

func (s *SwapMessage) Serialize() []byte {
	payloadLen := len(s.Payload)
	data := make([]byte, 2 + payloadLen)

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(s.Type))
	copy(data[:2], buf)
	copy(data[2:], s.Payload)

	return []byte(base64.StdEncoding.EncodeToString(data))
}

func DeserializeSwapMessage(data []byte) (*SwapMessage, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	messageType := binary.BigEndian.Uint16(decodedData[:2])

	return &SwapMessage{
		Type:    SwapMessageType(messageType),
		Payload: decodedData[2:],
	}, nil
}
