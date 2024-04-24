package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"

	ed25519 "filippo.io/edwards25519"
)

type Seed = []byte

type PrivateSpendKey struct {
	key *ed25519.Scalar
}

// View returns the private view key corresponding to the PrivateSpendKey.
func (k *PrivateSpendKey) View() (*PrivateViewKey, error) {
	h := Keccak256(k.key.Bytes())

	// We cannot utilize the `SetBytesWithClamping` function here,
	// which would perform the `sc_reduce32` computation on our behalf.
	//
	// The reason is that standard Monero wallets do not alter the first
	// and last byte during the view key calculation process.
	vkBytes := scReduce32(h)
	vk, err := ed25519.NewScalar().SetCanonicalBytes(vkBytes[:])
	if err != nil {
		return nil, err
	}

	return &PrivateViewKey{key: vk}, nil
}
func (sk *PrivateSpendKey) AsPrivateKeyPair() (*KeyPair, error) {
	vk, err := sk.View()
	if err != nil {
		return nil, err
	}
	newKeyPair := &KeyPair{sk: sk, vk: vk}

	return newKeyPair, nil
}

type PrivateViewKey struct {
	key *ed25519.Scalar
}

type KeyPair struct {
	sk *PrivateSpendKey
	vk *PrivateViewKey
}

func getSeed() (Seed, error) {
	// Generate a new random seed
	var seed Seed
	_, err := rand.Read(seed[:])
	if err != nil {
		return seed, err
	}

	// we hash the seed for compatibility w/ the ed25519 stdlib
	h := sha512.Sum512(seed[:])

	return h[:32], nil
}

func NewWallet() (*KeyPair, error) {
	seed, err := getSeed()
	if err != nil {
		return nil, fmt.Errorf("failed to get seed: %w", err)
	}

	s, err := ed25519.NewScalar().SetBytesWithClamping(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to set bytes: %w", err)
	}

	sk := &PrivateSpendKey{key: s}

	return sk.AsPrivateKeyPair()
}
