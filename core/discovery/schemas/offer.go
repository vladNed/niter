package schemas

type OfferMessage struct {
	Type     string `json:"type"`
	OfferID  string `json:"offerId"`
	OfferSDP string `json:"offerDescription"`
}
