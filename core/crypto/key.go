package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"github.com/c0mm4nd/go-ripemd"
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

func Hash160(data []byte) []byte {
	shaHasher := sha256.New()
	shaHasher.Write(data)
	shaHash := shaHasher.Sum(nil)

	ripemdHasher := ripemd.New160()
	ripemdHasher.Write(shaHash)
	return ripemdHasher.Sum(nil)
}
