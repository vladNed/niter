package discovery

import (
	"github.com/indexone/niter/core/discovery/schemas"
	"github.com/indexone/niter/core/p2p/protocol"
)

var Cache = NewOffersCache()

type OffersCache struct {
	Offers map[string]schemas.OfferMessage
}

func NewOffersCache() *OffersCache {
	return &OffersCache{Offers: make(map[string]schemas.OfferMessage)}
}

func (oc *OffersCache) AddOffer(offer schemas.OfferMessage) {
	oc.Offers[offer.OfferID] = offer
}

func (oc *OffersCache) AllOffers() []interface{} {
	offers := make([]interface{}, 0)
	for offerId := range oc.Offers {
		data := map[string]interface{}{
			"id": offerId,
			"sendingAmount": protocol.ConvertToFloat(
				oc.Offers[offerId].OfferDetails.SendingAmount,
				oc.Offers[offerId].OfferDetails.SendingCurrency,
			),
			"sendingCurrency": oc.Offers[offerId].OfferDetails.SendingCurrency,
			"receivingAmount": protocol.ConvertToFloat(
				oc.Offers[offerId].OfferDetails.ReceivingAmount,
				oc.Offers[offerId].OfferDetails.ReceivingCurrency,
			),
			"receivingCurrency": oc.Offers[offerId].OfferDetails.ReceivingCurrency,
			"swapCreator": oc.Offers[offerId].OfferDetails.SwapCreator,
		}
		offers = append(offers, data)
	}

	return offers
}

func (oc *OffersCache) GetOffer(offerID string) (schemas.OfferMessage, bool) {
	offer, ok := oc.Offers[offerID]
	return offer, ok
}

func (oc *OffersCache) RemoveOffer(offerID string) {
	delete(oc.Offers, offerID)
}
