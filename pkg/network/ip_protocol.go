package network

import "net"

/*

TUN - IPv4 Packet:
  +---------------------------------------------------------------------------------------------------------------+
  |       | Octet |           0           |           1           |           2           |           3           |
  | Octet |  Bit  |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|
  +---------------------------------------------------------------------------------------------------------------+
  |   0   |   0   |  Version  |    IHL    |      DSCP       | ECN |                 Total  Length                 |
  +---------------------------------------------------------------------------------------------------------------+
  |   4   |  32   |                Identification                 | Flags  |           Fragment Offset            |
  +---------------------------------------------------------------------------------------------------------------+
  |   8   |  64   |     Time To Live      |       Protocol        |                Header Checksum                |
  +---------------------------------------------------------------------------------------------------------------+
  |  12   |  96   |                                       Source IP Address                                       |
  +---------------------------------------------------------------------------------------------------------------+
  |  16   |  128  |                                    Destination IP Address                                     |
  +---------------------------------------------------------------------------------------------------------------+
  |  20   |  160  |                                     Options (if IHL > 5)                                      |
  +---------------------------------------------------------------------------------------------------------------+
  |  24   |  192  |                                                                                               |
  |  30   |  224  |                                            Payload                                            |
  |  ...  |  ...  |                                                                                               |
  +---------------------------------------------------------------------------------------------------------------+

*/

type IPProtocol byte

var IPProtocols = map[byte]string{
	0x00: "POPTHO",
	0x01: "ICMP",
	0x02: "IGMP",
	0x03: "GGP",
	0x04: "IPv4Encapsulation",
	0x05: "ST",
	0x06: "TCP",
	0x07: "CBT",
	0x08: "EGP",
	0x09: "IGP",
	0x0A: "BBN_RCC_MON",
	0x0B: "NVP_II",
	0x0C: "PUP",
	0x0D: "ARGUS",
	0x0E: "EMCON",
	0x0F: "XNET",
	0x10: "CHAOS",
	0x11: "UDP",
	0x12: "MUX",
	0x13: "DCN_MEAS",
	0x14: "HMP",
	0x15: "PRM",
	0x16: "XNS_IDP",
	0x17: "TRUNK_1",
	0x18: "TRUNK_2",
	0x19: "LEAF_1",
	0x1A: "LEAF_2",
	0x1B: "RDP",
	0x1C: "IRTP",
	0x1D: "ISO_TP4",
	0x1E: "NETBLT",
	0x1F: "MFE_NSP",
	0x20: "MERIT_INP",
	0x21: "DCCP",
	0x22: "ThirdPC",
	0x23: "IDPR",
	0x24: "XTP",
	0x25: "DDP",
	0x26: "IDPR_CMTP",
	0x27: "TPxx",
	0x28: "IL",
	0x29: "IPv6Encapsulation",
	0x2A: "SDRP",
	0x2B: "IPv6_Route",
	0x2C: "IPv6_Frag",
	0x2D: "IDRP",
	0x2E: "RSVP",
	0x2F: "GRE",
	0x30: "MHRP",
	0x31: "BNA",
	0x32: "ESP",
	0x33: "AH",
	0x34: "I_NLSP",
	0x35: "SWIPE",
	0x36: "NARP",
	0x37: "MOBILE",
	0x38: "TLSP",
	0x39: "SKIP",
	0x3A: "IPv6_ICMP",
	0x3B: "IPv6_NoNxt",
	0x3C: "IPv6_Opts",
	0x3E: "CFTP",
	0x40: "SAT_EXPAK",
	0x41: "KRYPTOLAN",
	0x42: "RVD",
	0x43: "IPPC",
	0x45: "SAT_MON",
	0x46: "VISA",
	0x47: "IPCV",
	0x48: "CPNX",
	0x49: "CPHB",
	0x4A: "WSN",
	0x4B: "PVP",
	0x4C: "BR_SAT_MON",
	0x4D: "SUN_ND",
	0x4E: "WB_MON",
	0x4F: "WB_EXPAK",
	0x50: "ISO_IP",
	0x51: "VMTP",
	0x52: "SECURE_VMTP",
	0x53: "VINES",
	0x54: "IPTM",
	0x55: "NSFNET_IGP",
	0x56: "DGP",
	0x57: "TCF",
	0x58: "EIGRP",
	0x59: "OSPF",
	0x5A: "Sprite_RPC",
	0x5B: "LARP",
	0x5C: "MTP",
	0x5D: "AX_25",
	0x5E: "IPIP",
	0x5F: "MICP",
	0x60: "SCC_SP",
	0x61: "ETHERIP",
	0x62: "ENCAP",
	0x64: "GMTP",
	0x65: "IFMP",
	0x66: "PNNI",
	0x67: "PIM",
	0x68: "ARIS",
	0x69: "SCPS",
	0x6A: "QNX",
	0x6B: "A_N",
	0x6C: "IPComp",
	0x6D: "SNP",
	0x6E: "Compaq_Peer",
	0x6F: "IPX_in_IP",
	0x70: "VRRP",
	0x71: "PGM",
	0x73: "L2TP",
	0x74: "DDX",
	0x75: "IATP",
	0x76: "STP",
	0x77: "SRP",
	0x78: "UTI",
	0x79: "SMP",
	0x7A: "SM",
	0x7B: "PTP",
	0x7D: "FIRE",
	0x7E: "CRTP",
	0x7F: "CRUDP",
	0x80: "SSCOPMCE",
	0x81: "IPLT",
	0x82: "SPS",
	0x83: "PIPE",
	0x84: "SCTP",
	0x85: "FC",
	0x8A: "manet",
	0x8B: "HIP",
	0x8C: "Shim6",
}

func IsIPv4(packet []byte) bool {
	return 4 == (packet[0] >> 4)
}

func IsIPv6(packet []byte) bool {
	return 6 == (packet[0] >> 4)
}

func IPv4DSCP(packet []byte) byte {
	return packet[1] >> 2
}

func IPv4ECN(packet []byte) byte {
	return packet[1] & 0x03
}

func IPv4Identification(packet []byte) [2]byte {
	return [2]byte{packet[4], packet[5]}
}

func IPv4TTL(packet []byte) byte {
	return packet[8]
}

func IPv4Protocol(packet []byte) byte {
	return packet[9]
}

func IPv4Source(packet []byte) net.IP {
	return net.IPv4(packet[12], packet[13], packet[14], packet[15])
}

func SetIPv4Source(packet []byte, source net.IP) {
	copy(packet[12:16], source.To4())
}

func IPv4Destination(packet []byte) net.IP {
	return net.IPv4(packet[16], packet[17], packet[18], packet[19])
}

func SetIPv4Destination(packet []byte, dest net.IP) {
	copy(packet[16:20], dest.To4())
}

func IPv4Payload(packet []byte) []byte {
	ihl := packet[0] & 0x0F
	return packet[ihl*4:]
}

// For TCP/UDP
func IPv4SourcePort(packet []byte) uint16 {
	payload := IPv4Payload(packet)
	return (uint16(payload[0]) << 8) | uint16(payload[1])
}

func IPv4DestinationPort(packet []byte) uint16 {
	payload := IPv4Payload(packet)
	return (uint16(payload[2]) << 8) | uint16(payload[3])
}

func SetIPv4SourcePort(packet []byte, port uint16) {
	payload := IPv4Payload(packet)
	payload[0] = byte(port >> 8)
	payload[1] = byte(port & 0xFF)
}

func SetIPv4DestinationPort(packet []byte, port uint16) {
	payload := IPv4Payload(packet)
	payload[2] = byte(port >> 8)
	payload[3] = byte(port & 0xFF)
}
