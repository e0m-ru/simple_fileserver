package config

import (
	"testing"
)

func TestGetAllLocalIPs(t *testing.T) {
	ips, err := GetAllLocalIPs()
	if err != nil {
		t.Fatalf("GetAllLocalIPs() error = %v", err)
	}

	if len(ips) == 0 {
		t.Logf("No local IPs found (might be expected in some environments)")
	}

	// Check that returned IPs are not loopback
	for _, ip := range ips {
		if ip == "127.0.0.1" || ip == "::1" {
			t.Errorf("Loopback IP should not be returned: %s", ip)
		}
	}

	t.Logf("Found %d local IPs", len(ips))
}
