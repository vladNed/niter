package schemas

import (
	"encoding/json"
	"errors"
)

var (
	Offer  = "offer"
	Answer = "answer"
)

// Message is the interface that all message types must implement.
type Message interface{}

// ReceivedMessage is the base message type that all messages must contain.
type ReceivedMessage struct {
	Type string `json:"type"`
}

// unmarshal unmarshals the payload into the appropriate message type.
func unmarshal(msg *ReceivedMessage, payload *[]byte) (Message, error) {
	switch msg.Type {
	case Offer:
		var req OfferMessage
		return &req, json.Unmarshal(*payload, &req)
	case Answer:
		var req AnswerMessage
		return &req, json.Unmarshal(*payload, &req)
	default:
		return nil, errors.New("unknown message type")
	}

}

// ParseReceivedMessage parses a received message and returns the appropriate message type
// based on the message type field.
func ParseReceivedMessage(payload []byte) (Message, error) {
	var response ReceivedMessage
	err := json.Unmarshal(payload, &response)
	if err != nil {
		return nil, err
	}
	msg, err := unmarshal(&response, &payload)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
