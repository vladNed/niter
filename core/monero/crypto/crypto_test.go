package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashing(t *testing.T) {
	// Test Keccak256
	data := []byte("hello")
	hash := Keccak256(data)
	assert.Equal(t, 32, len(hash), "Keccak256 failed")

	// Test scReduce32
	reduced := scReduce32(hash)
	assert.Equal(t, 32, len(reduced), "scReduce32 failed")
}