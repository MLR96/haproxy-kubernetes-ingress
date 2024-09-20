package zone

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var controllerIP string
var controllerZone string
var zonesInfo string
var zonesCrossEnabled bool = true

func init() {
	controllerIP = os.Getenv("POD_IP")
	zonesInfo = os.Getenv("ZONES_INFO")
	zonesCrossEnabled = !(strings.ToLower(os.Getenv("ZONES_CROSS_TRAFFIC_ENABLED")) == "false")

	if controllerIP == "" || zonesInfo == "" {
		zonesCrossEnabled = true
		return
	}

	controllerZone = getZoneFromIP(controllerIP, zonesInfo)
}

// This function is used when haproxy configure a server in backend
func IsBackupEnabledForThisIP(address string) bool {
	if zonesCrossEnabled || address == "" || address == "127.0.0.1" {
		return false
	}
	srvZone := getZoneFromIP(address, zonesInfo)
	return srvZone != controllerZone
}

// Private; Auxiliary functions
func getZoneFromIP(address string, zoneSubnets string) string {
	input := strings.ReplaceAll(zoneSubnets, " ", "")
	subnetGroups := strings.Split(input, ";")

	for _, group := range subnetGroups {
		parts := strings.Split(group, ":")
		if len(parts) != 2 || parts[1] == "" {
			return "unknown"
		}

		zoneName := parts[0]
		subnets := strings.Split(parts[1], ",")

		for _, subnet := range subnets {
			if ok, _ := isIPInSubnet(address, subnet); ok {
				return zoneName
			}
		}
	}

	return "unknown"
}

func isIPInSubnet(ipStr, subnetStr string) (bool, error) {
	// Parseamos la IP
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("IP inválida: %s", ipStr)
	}

	// Parseamos la subnet (dirección IP + máscara)
	_, subnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return false, fmt.Errorf("subred inválida: %s", subnetStr)
	}

	// Comprobamos si la IP está dentro de la subred
	return subnet.Contains(ip), nil
}
