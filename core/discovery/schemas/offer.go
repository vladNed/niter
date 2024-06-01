package schemas

type OfferMessage struct {
	Type         string       `json:"type"`
	OfferID      string       `json:"offerId"`
	OfferSDP     string       `json:"offerDescription"`
	OfferDetails OfferDetails `json:"offerDetails"`
}

type OfferDetails struct {
	SwapCreator       string `json:"swapCreator"`
	SendingAmount     string `json:"sendingAmount"`
	SendingCurrency   string `json:"sendingCurrency"`
	ReceivingAmount   string `json:"receivingAmount"`
	ReceivingCurrency string `json:"receivingCurrency"`
	CreatedAt         string `json:"createdAt"`
	ExpiresAt         string `json:"expiresAt"`
}
