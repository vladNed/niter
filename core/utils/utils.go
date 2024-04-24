package utils

import (
	"encoding/json"
	"encoding/base64"

	"github.com/pion/webrtc/v4"
)

// Encodes to base64 an SDP object to be sent to the signalling server
// as url encoded form data
func EncodeSDP(decodedSDP *webrtc.SessionDescription) (string, error) {
	sdpJson, err := json.Marshal(decodedSDP)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sdpJson), nil
}

// Decodes a base64 encoded SDP object from the signalling server
func DecodeSDP(encodedSDP string) (*webrtc.SessionDescription, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedSDP)
	if err != nil {
		return nil, err
	}

	var sdp webrtc.SessionDescription
	err = json.Unmarshal(decoded, &sdp)
	if err != nil {
		return nil, err
	}

	return &sdp, nil
}
