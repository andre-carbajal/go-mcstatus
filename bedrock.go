package gomcstat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var unconnectedPingData = []byte{
	0x01,                                           // Packet ID
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Time
	0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78, // Offline Message Data ID
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Client GUID
}

type BedrockServer struct {
	Host string
	Port uint16
}

func NewBedrockServer(address string) (*BedrockServer, error) {
	// Bedrock doesn't use SRV
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		if strings.Contains(address, ":") {
			return nil, err
		}
		host = address
		portStr = "19132"
	}
	port, _ := strconv.ParseUint(portStr, 10, 16)
	return &BedrockServer{Host: host, Port: uint16(port)}, nil
}

func (s *BedrockServer) Status() (StatusResponse, error) {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err := net.DialTimeout("udp", addr, 3*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	start := time.Now()
	if _, err := conn.Write(unconnectedPingData); err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	latency := time.Since(start).Milliseconds()

	if n < 35 || buf[0] != 0x1C {
		return nil, errors.New("invalid bedrock response")
	}

	serverIDLength := binary.BigEndian.Uint16(buf[33:35])
	if n < 35+int(serverIDLength) {
		return nil, errors.New("bedrock response too short")
	}

	dataStr := string(buf[35 : 35+serverIDLength])
	parts := strings.Split(dataStr, ";")

	if len(parts) < 6 {
		return nil, errors.New("not enough data in bedrock response")
	}

	protocol, _ := strconv.Atoi(parts[2])
	online, _ := strconv.Atoi(parts[4])
	max, _ := strconv.Atoi(parts[5])

	resp := &BedrockStatusResponse{
		ServerID: parts[0],
		MOTD:     parts[1],
		Protocol: protocol,
		Version:  parts[3],
		Online:   online,
		Max:      max,
		Latency:  latency,
	}

	if len(parts) > 7 {
		resp.MapName = parts[7]
	}
	if len(parts) > 8 {
		resp.Gamemode = parts[8]
	}

	return resp, nil
}

func (s *BedrockServer) Ping() (int64, error) {
	status, err := s.Status()
	if err != nil {
		return 0, err
	}
	return status.GetLatency(), nil
}
