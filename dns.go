package mcstatus

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ResolveAddress(address string, defaultPort uint16) (string, uint16, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		// No port specified, or invalid format
		if strings.Contains(address, ":") {
			return "", 0, fmt.Errorf("invalid address format: %w", err)
		}
		host = address
		// Try SRV lookup
		_, srvs, err := net.LookupSRV("minecraft", "tcp", host)
		if err == nil && len(srvs) > 0 {
			// Remove trailing dot from target
			target := strings.TrimSuffix(srvs[0].Target, ".")
			return target, srvs[0].Port, nil
		}
		return host, defaultPort, nil
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %w", err)
	}

	return host, uint16(port), nil
}
