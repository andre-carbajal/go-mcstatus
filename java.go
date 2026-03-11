package mcstatus

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type JavaServer struct {
	Host string
	Port uint16
}

func NewJavaServer(address string) (*JavaServer, error) {
	host, port, err := ResolveAddress(address, 25565)
	if err != nil {
		return nil, err
	}
	return &JavaServer{Host: host, Port: port}, nil
}

func (s *JavaServer) Status() (StatusResponse, error) {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	// Handshake
	var handshake PacketBuffer
	handshake.WriteVarInt(0x00) // Packet ID
	handshake.WriteVarInt(47)   // Protocol version
	handshake.WriteString(s.Host)
	handshake.WriteUShort(s.Port)
	handshake.WriteVarInt(1) // Next state: Status

	if err := writePacket(conn, &handshake); err != nil {
		return nil, err
	}

	// Request
	var request PacketBuffer
	request.WriteVarInt(0x00)
	start := time.Now()
	if err := writePacket(conn, &request); err != nil {
		return nil, err
	}

	// Response
	reader := bufio.NewReader(conn)
	length, err := ReadVarInt(reader)
	if err != nil {
		return nil, err
	}

	packetBytes := make([]byte, length)
	if _, err := io.ReadFull(reader, packetBytes); err != nil {
		return nil, err
	}

	packetReader := bytes.NewReader(packetBytes)
	packetID, err := ReadVarInt(packetReader)
	if err != nil {
		return nil, err
	}

	if packetID != 0x00 {
		return nil, errors.New("invalid status response packet")
	}

	jsonStr, err := ReadString(packetReader)
	if err != nil {
		return nil, err
	}

	latency := time.Since(start).Milliseconds()

	var resp JavaStatusResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, err
	}
	resp.Latency = latency

	return &resp, nil
}

func (s *JavaServer) Ping() (int64, error) {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	// Handshake
	var handshake PacketBuffer
	handshake.WriteVarInt(0x00) // Packet ID
	handshake.WriteVarInt(47)   // Protocol version
	handshake.WriteString(s.Host)
	handshake.WriteUShort(s.Port)
	handshake.WriteVarInt(1) // Next state: Status

	if err := writePacket(conn, &handshake); err != nil {
		return 0, err
	}

	// Request
	var request PacketBuffer
	request.WriteVarInt(0x00)
	if err := writePacket(conn, &request); err != nil {
		return 0, err
	}

	// Response (Read and discard)
	reader := bufio.NewReader(conn)
	length, err := ReadVarInt(reader)
	if err != nil {
		return 0, err
	}
	packetBytes := make([]byte, length)
	if _, err := io.ReadFull(reader, packetBytes); err != nil {
		return 0, err
	}

	// Ping
	var ping PacketBuffer
	ping.WriteVarInt(0x01)
	pingTime := time.Now().UnixNano()
	binary.Write(&ping, binary.BigEndian, pingTime)

	start := time.Now()
	if err := writePacket(conn, &ping); err != nil {
		return 0, err
	}

	// Pong
	reader = bufio.NewReader(conn)
	length, err = ReadVarInt(reader)
	if err != nil {
		return 0, err
	}

	packetBytes = make([]byte, length)
	if _, err = io.ReadFull(reader, packetBytes); err != nil {
		return 0, err
	}

	packetReader := bytes.NewReader(packetBytes)
	packetID, err := ReadVarInt(packetReader)
	if err != nil || packetID != 0x01 {
		return 0, errors.New("invalid ping response packet")
	}

	latency := time.Since(start).Milliseconds()
	return latency, nil
}

func writePacket(conn net.Conn, packet *PacketBuffer) error {
	var wrapper PacketBuffer
	wrapper.WriteVarInt(packet.Len())
	_, err := conn.Write(append(wrapper.Bytes(), packet.Bytes()...))
	return err
}
