package discovery

import (
	"reflect"

	"github.com/indexone/niter/core/discovery/schemas"
)

var Cache = NewOffersCache()

type CacheData = map[string]interface{}

type OffersCache struct {
	Offers map[string]CacheData
}

func NewOffersCache() *OffersCache {
	return &OffersCache{Offers: make(map[string]CacheData)}
}

func (oc *OffersCache) AddOffer(offer schemas.OfferMessage) {
	data := oc.structToMap(offer)
	oc.Offers[offer.OfferID] = data
}

func (oc *OffersCache) AllOffers() []interface{} {
	offers := make([]interface{}, 0)
	for offerId := range oc.Offers {
		data := map[string]interface{}{
			"id": offerId,
		}
		offers = append(offers, data)
	}

	return offers
}

func (oc *OffersCache) GetOffer(offerID string) (CacheData, bool) {
	offer, ok := oc.Offers[offerID]
	return offer, ok
}

func (oc *OffersCache) RemoveOffer(offerID string) {
	delete(oc.Offers, offerID)
}

func (oc *OffersCache) structToMap(data interface{}) CacheData {
	result := make(map[string]interface{})

	// Using reflection to get struct fields and values
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()
		result[fieldName] = fieldValue
	}

	return result
}
