package discovery

import "testing"

func TestNewMDNSClient(t *testing.T) {
	_, err := NewMDNSClient()
	if err != nil {
		t.Fatalf(err.Error())
	}
}
