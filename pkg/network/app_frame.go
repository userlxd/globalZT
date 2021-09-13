package network

import (
	"bytes"
	"encoding/binary"
	"globalZT/pkg/tunnel"
	"net"
	"strconv"
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

func Parse(packet []byte) (string, *tunnel.OfficeReq) {
	var code = bytes.Buffer{}
	var req = &tunnel.OfficeReq{}

	// IP Protocol
	code.WriteString(IPProtocols[packet[9]])
	code.WriteString(":")
	if packet[9] == 0x06 {
		req.Protocol = 1
	}
	if packet[9] == 0x11 {
		req.Protocol = 2
	}

	// Source IP
	code.WriteString(net.IPv4(packet[12], packet[13], packet[14], packet[15]).String())
	req.Sip = binary.BigEndian.Uint32(packet[12:16])

	code.WriteString(">")

	// Destination IP
	code.WriteString(net.IPv4(packet[16], packet[17], packet[18], packet[19]).String())
	req.Sip = binary.BigEndian.Uint32(packet[16:20])
	code.WriteString(":")

	// IP Payload
	ihl := packet[0] & 0x0F
	payload := packet[ihl*4:]

	// Transport Layer Port
	req.Port = binary.BigEndian.Uint32(payload[0:4])
	code.WriteString(strconv.Itoa(int((uint16(payload[0]) << 8) | uint16(payload[1]))))
	code.WriteString(">")
	code.WriteString(strconv.Itoa(int((uint16(payload[2]) << 8) | uint16(payload[3]))))

	// Data
	req.Data = packet

	// Host
	req.Host = "a.b.c"

	return code.String(), req
}
