package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	// Test NewWallet
	kp, err := NewWallet()
	assert.NoError(t, err, "NewWallet failed")

	// Check if the returned KeyPair is not nil
	assert.NotNil(t, kp, "Returned KeyPair is nil")

	// Check if the PrivateSpendKey and PrivateViewKey are not nil
	assert.NotNil(t, kp.sk, "PrivateSpendKey is nil")
	assert.NotNil(t, kp.vk, "PrivateViewKey is nil")
}
