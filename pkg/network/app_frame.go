package network

import (
	"encoding/binary"
	"globalZT/pkg/tunnel"
	"net"
)

type ReqFrame struct {
	Code     string
	Sip      net.IP
	Dip      net.IP
	Sport    uint16
	Dport    uint16
	Protocol string
	Packet   []byte
}

func Parse(packet []byte) *tunnel.In {
	var req = &tunnel.In{}

	// IP Protocol
	if packet[9] == 0x06 {
		req.Protocol = 1
	}
	if packet[9] == 0x11 {
		req.Protocol = 2
	}

	// Source IP
	req.Sip = binary.BigEndian.Uint32(packet[12:16])

	// Destination IP
	req.Sip = binary.BigEndian.Uint32(packet[16:20])

	// IP Payload
	ihl := packet[0] & 0x0F
	payload := packet[ihl*4:]

	// Transport Layer Port
	req.Port = binary.BigEndian.Uint32(payload[0:4])
	// Data
	req.Data = packet

	// Host
	req.Host = "a.b.c"

	return req
}
