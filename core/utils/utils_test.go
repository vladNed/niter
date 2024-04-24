package utils

import (
	"testing"

	"github.com/pion/webrtc/v4"
)

func TestEncodeSDP(t *testing.T) {
	sdp := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP: "v=0\r\no=- 0 0 IN IP4",
	}

	encoded, err := EncodeSDP(&sdp)
	if err != nil {
		t.Errorf("Error encoding SDP: %v", err)
	}

	if encoded == "" {
		t.Errorf("Empty encoded SDP")
	}
}

func TestDecodeSDP(t *testing.T) {
	sdp := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP: "dummy-data",
	}

	encoded, err := EncodeSDP(&sdp)
	if err != nil {
		t.Errorf("Error encoding SDP: %v", err)
	}

	decoded, err := DecodeSDP(encoded)
	if err != nil {
		t.Errorf("Error decoding SDP: %v", err)
	}

	if decoded.Type != sdp.Type && decoded.SDP != sdp.SDP{
		t.Errorf("Decoded SDP type does not match")
	}

}

func TestDecodeSDPErr(t *testing.T) {
	_, err := DecodeSDP("dummy-data")
	if err == nil {
		t.Errorf("Expected error decoding SDP")
	}
}