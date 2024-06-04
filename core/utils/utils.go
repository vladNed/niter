package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

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

func Hash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)

	return hex.EncodeToString(hash.Sum(nil))
}

func GetTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}


func ToIntArray(data []byte) []int {
	intData := make([]int, len(data))
	for i, b := range data {
		intData[i] = int(b)
	}

	return intData
}

func ToByteArray(data []int) []byte {
	byteData := make([]byte, len(data))
	for i, b := range data {
		byteData[i] = byte(b)
	}

	return byteData
}

func GenerateSeed() ([]byte, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}


func ToBigInt(n string) *big.Int {
	val := new(big.Int)
	val.SetString(n, 10)

	return val
}