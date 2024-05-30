package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

)

type NetworkKey struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func (k *NetworkKey) Commitment() string {
	hasher := sha256.New()
	hasher.Write(k.PublicKey)
	firstHash := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(firstHash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateKey() (*NetworkKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &NetworkKey{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}