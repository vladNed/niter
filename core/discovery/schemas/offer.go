package schemas

import (
	"github.com/pion/webrtc/v4"
)

type OfferMessage struct {
	Type             string                    `json:"type"`
	OfferID          string                    `json:"offerId"`
	OfferDescription *webrtc.SessionDescription `json:"offerDescription"`
}
