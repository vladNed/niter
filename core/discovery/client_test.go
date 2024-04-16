package discovery

import (
	"testing"
)

func TestErrorOnWSClientCreate(t *testing.T) {
	_, err := NewWSClient()
	if err == nil {
		t.Errorf("No error on WSClient creation")
	}
}