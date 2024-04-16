package schemas

type AnswerMessage struct {
	Type      string `json:"type"`
	OfferID   string `json:"offerId"`
	AnswerSDP string `json:"sdp"`
}
