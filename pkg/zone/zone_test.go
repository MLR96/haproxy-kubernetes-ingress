package zone

import (
	"testing"
)

func TestGetZoneFromIp(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		zoneSubnets string
		expected    string
	}{
		{
			name:        "IP belongs to eu-south-2a",
			address:     "192.168.2.1",
			zoneSubnets: "eu-south-2a: 192.168.1.0/24, 192.168.2.0/24; eu-south-2b: 192.168.3.0/24",
			expected:    "eu-south-2a",
		},
		{
			name:        "IP belongs to eu-south-2b",
			address:     "192.168.3.10",
			zoneSubnets: "eu-south-2a: 192.168.1.0/24, 192.168.2.0/24; eu-south-2b: 192.168.3.0/24",
			expected:    "eu-south-2b",
		},
		{
			name:        "IP not in any subnet",
			address:     "10.0.0.1",
			zoneSubnets: "eu-south-2a: 192.168.1.0/24, 192.168.2.0/24; eu-south-2b: 192.168.3.0/24",
			expected:    "unknown",
		},
		{
			name:        "Malformed input - no subnet",
			address:     "192.168.2.1",
			zoneSubnets: "eu-south-2a:;eu-south-2b:192.168.2.0/24",
			expected:    "unknown",
		},
		{
			name:        "Malformed input - no subnet",
			address:     "192.168.2.1",
			zoneSubnets: "eu-south-2a;eu-south-2b:192.168.2.0/24",
			expected:    "unknown",
		},
		{
			name:        "IP belongs to overlapping subnets, first match",
			address:     "192.167.4.10",
			zoneSubnets: "eu-south-2a:192.167.4.0/24;eu-south-2b:192.167.4.0/16",
			expected:    "eu-south-2a",
		},
		{
			name:        "Empty Data",
			address:     "192.167.4.10",
			zoneSubnets: "   ",
			expected:    "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getZoneFromIP(tt.address, tt.zoneSubnets)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
