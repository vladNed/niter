package schemas

import "github.com/pion/webrtc/v4"

type AnswerMessage struct {
	Type              string                     `json:"type"`
	OfferID           string                     `json:"offerId"`
	AnswerDescription *webrtc.SessionDescription `json:"answerDescription"`
}
