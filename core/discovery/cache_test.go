package discovery

import (
	"testing"

	"github.com/indexone/niter/core/discovery/schemas"
)

func TestOffersCache(t *testing.T) {
	// Create a new OffersCache
	cache := NewOffersCache()

	// Test AddOffer
	offer := schemas.OfferMessage{OfferID: "1" /* other fields */}
	cache.AddOffer(offer)

	// Test GetOffer
	data, ok := cache.GetOffer("1")
	if !ok || data["OfferID"] != "1" {
		t.Fatalf("GetOffer failed: offer not found or data does not match")
	}

	// Test AllOffers
	offers := cache.AllOffers()
	if len(offers) != 1 || offers[0] != "1" {
		t.Fatalf("AllOffers failed: offer not found or data does not match")
	}

	// Test RemoveOffer
	cache.RemoveOffer("1")
	_, ok = cache.GetOffer("1")
	if ok {
		t.Fatalf("RemoveOffer failed: offer still found after removal")
	}
}
