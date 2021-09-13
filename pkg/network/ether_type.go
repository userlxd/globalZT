package network

import "net"

/*

TAP - MAC Frame:
   No Tagging
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   |EType| Payload
  +-----------------------------------------------------------------------------

   Single-Tagged -- Octets [12,13] == {0x81, 0x00}
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   |    Tag    | Payload
  +-----------------------------------------------------------------------------

   Double-Tagged -- Octets [12,13] == {0x88, 0xA8}
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   | Outer Tag | Inner Tag | Payload
  +-----------------------------------------------------------------------------

*/

type Ethertype [2]byte

var (
	IPv4                = Ethertype{0x08, 0x00}
	ARP                 = Ethertype{0x08, 0x06}
	WakeOnLAN           = Ethertype{0x08, 0x42}
	TRILL               = Ethertype{0x22, 0xF3}
	DECnetPhase4        = Ethertype{0x60, 0x03}
	RARP                = Ethertype{0x80, 0x35}
	AppleTalk           = Ethertype{0x80, 0x9B}
	AARP                = Ethertype{0x80, 0xF3}
	IPX1                = Ethertype{0x81, 0x37}
	IPX2                = Ethertype{0x81, 0x38}
	QNXQnet             = Ethertype{0x82, 0x04}
	IPv6                = Ethertype{0x86, 0xDD}
	EthernetFlowControl = Ethertype{0x88, 0x08}
	IEEE802_3           = Ethertype{0x88, 0x09}
	CobraNet            = Ethertype{0x88, 0x19}
	MPLSUnicast         = Ethertype{0x88, 0x47}
	MPLSMulticast       = Ethertype{0x88, 0x48}
	PPPoEDiscovery      = Ethertype{0x88, 0x63}
	PPPoESession        = Ethertype{0x88, 0x64}
	JumboFrames         = Ethertype{0x88, 0x70}
	HomePlug1_0MME      = Ethertype{0x88, 0x7B}
	IEEE802_1X          = Ethertype{0x88, 0x8E}
	PROFINET            = Ethertype{0x88, 0x92}
	HyperSCSI           = Ethertype{0x88, 0x9A}
	AoE                 = Ethertype{0x88, 0xA2}
	EtherCAT            = Ethertype{0x88, 0xA4}
	EthernetPowerlink   = Ethertype{0x88, 0xAB}
	LLDP                = Ethertype{0x88, 0xCC}
	SERCOS3             = Ethertype{0x88, 0xCD}
	HomePlugAVMME       = Ethertype{0x88, 0xE1}
	MRP                 = Ethertype{0x88, 0xE3}
	IEEE802_1AE         = Ethertype{0x88, 0xE5}
	IEEE1588            = Ethertype{0x88, 0xF7}
	IEEE802_1ag         = Ethertype{0x89, 0x02}
	FCoE                = Ethertype{0x89, 0x06}
	FCoEInit            = Ethertype{0x89, 0x14}
	RoCE                = Ethertype{0x89, 0x15}
	CTP                 = Ethertype{0x90, 0x00}
	VeritasLLT          = Ethertype{0xCA, 0xFE}
)

type Tagging int

// Indicating whether/how a MAC frame is tagged. The value is number of bytes taken by tagging.
const (
	NotTagged    Tagging = 0
	Tagged       Tagging = 4
	DoubleTagged Tagging = 8
)

func MACDestination(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[:6])
}

func MACSource(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[6:12])
}

func MACTagging(macFrame []byte) Tagging {
	if macFrame[12] == 0x81 && macFrame[13] == 0x00 {
		return Tagged
	} else if macFrame[12] == 0x88 && macFrame[13] == 0xa8 {
		return DoubleTagged
	}
	return NotTagged
}

func MACEthertype(macFrame []byte) Ethertype {
	ethertypePos := 12 + MACTagging(macFrame)
	return Ethertype{macFrame[ethertypePos], macFrame[ethertypePos+1]}
}

func MACPayload(macFrame []byte) []byte {
	return macFrame[12+MACTagging(macFrame)+2:]
}

func IsBroadcast(addr net.HardwareAddr) bool {
	return addr[0] == 0xff && addr[1] == 0xff && addr[2] == 0xff && addr[3] == 0xff && addr[4] == 0xff && addr[5] == 0xff
}

func IsIPv4Multicast(addr net.HardwareAddr) bool {
	return addr[0] == 0x01 && addr[1] == 0x00 && addr[2] == 0x5e
}
