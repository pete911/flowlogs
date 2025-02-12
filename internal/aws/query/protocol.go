package query

import (
	"fmt"
	"strconv"
	"strings"
)

type protocol struct {
	number   int
	keyword  string
	protocol string
}

func protocolFromKeywordToNumber(in string) int {
	for k, v := range protocolByNumber {
		if strings.ToLower(v.keyword) == strings.ToLower(in) {
			return k
		}
	}
	return -1
}

func ProtocolFromNumberToKeyword(in string) string {
	if v := toProtocol(in).keyword; v != "" {
		return v
	}
	return in
}

func toProtocol(in string) protocol {
	strVal := fmt.Sprintf("%s", in)
	if strVal == "-" || strVal == "" {
		return protocol{}
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return protocol{}
	}

	if intVal > 255 {
		return protocol{}
	}

	if intVal == 255 {
		return protocol{
			number:   intVal,
			keyword:  "Reserved",
			protocol: "",
		}
	}

	if intVal == 254 || intVal == 253 {
		return protocol{
			number:   intVal,
			keyword:  "",
			protocol: "Use for experimentation and testing",
		}
	}

	if intVal > 145 {
		return protocol{
			number:   intVal,
			keyword:  "",
			protocol: "Unassigned",
		}
	}

	if v, ok := protocolByNumber[intVal]; ok {
		return v
	}
	return protocol{}
}

// --- protocols ---

// ProtocolByNumber protocols mapped by their numbers - https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
var protocolByNumber = map[int]protocol{
	0: {
		number:   0,
		keyword:  "HOPOPT",
		protocol: "IPv6 Hop-by-Hop Option",
	},
	1: {
		number:   1,
		keyword:  "ICMP",
		protocol: "Internet Control Message",
	},
	2: {
		number:   2,
		keyword:  "IGMP",
		protocol: "Internet Group Management",
	},
	3: {
		number:   3,
		keyword:  "GGP",
		protocol: "Gateway-to-Gateway",
	},
	4: {
		number:   4,
		keyword:  "IPv4",
		protocol: "IPv4 encapsulation",
	},
	5: {
		number:   5,
		keyword:  "ST",
		protocol: "Stream",
	},
	6: {
		number:   6,
		keyword:  "TCP",
		protocol: "Transmission Control",
	},
	7: {
		number:   7,
		keyword:  "CBT",
		protocol: "CBT",
	},
	8: {
		number:   8,
		keyword:  "EGP",
		protocol: "Exterior Gateway ProtocolData",
	},
	9: {
		number:   9,
		keyword:  "IGP",
		protocol: "any private interior gateway (used by Cisco for their IGRP)",
	},
	10: {
		number:   10,
		keyword:  "BBN-RCC-MON",
		protocol: "BBN RCC Monitoring",
	},
	11: {
		number:   11,
		keyword:  "NVP-II",
		protocol: "Network Voice ProtocolData",
	},
	12: {
		number:   12,
		keyword:  "PUP",
		protocol: "PUP",
	},
	13: {
		number:   13,
		keyword:  "ARGUS (deprecated)",
		protocol: "ARGUS",
	},
	14: {
		number:   14,
		keyword:  "EMCON",
		protocol: "EMCON",
	},
	15: {
		number:   15,
		keyword:  "XNET",
		protocol: "Cross Net Debugger",
	},
	16: {
		number:   16,
		keyword:  "CHAOS",
		protocol: "Chaos",
	},
	17: {
		number:   17,
		keyword:  "UDP",
		protocol: "User Datagram",
	},
	18: {
		number:   18,
		keyword:  "MUX",
		protocol: "Multiplexing",
	},
	19: {
		number:   19,
		keyword:  "DCN-MEAS",
		protocol: "DCN Measurement Subsystems",
	},
	20: {
		number:   20,
		keyword:  "HMP",
		protocol: "Host Monitoring",
	},
	21: {
		number:   21,
		keyword:  "PRM",
		protocol: "Packet Radio Measurement",
	},
	22: {
		number:   22,
		keyword:  "XNS-IDP",
		protocol: "XEROX NS IDP",
	},
	23: {
		number:   23,
		keyword:  "TRUNK-1",
		protocol: "Trunk-1",
	},
	24: {
		number:   24,
		keyword:  "TRUNK-2",
		protocol: "Trunk-2",
	},
	25: {
		number:   25,
		keyword:  "LEAF-1",
		protocol: "Leaf-1",
	},
	26: {
		number:   26,
		keyword:  "LEAF-2",
		protocol: "Leaf-2",
	},
	27: {
		number:   27,
		keyword:  "RDP",
		protocol: "Reliable Data ProtocolData",
	},
	28: {
		number:   28,
		keyword:  "IRTP",
		protocol: "Internet Reliable Transaction",
	},
	29: {
		number:   29,
		keyword:  "ISO-TP4",
		protocol: "ISO Transport ProtocolData Class 4",
	},
	30: {
		number:   30,
		keyword:  "NETBLT",
		protocol: "Bulk Data Transfer ProtocolData",
	},
	31: {
		number:   31,
		keyword:  "MFE-NSP",
		protocol: "MFE Network Services ProtocolData",
	},
	32: {
		number:   32,
		keyword:  "MERIT-INP",
		protocol: "MERIT Internodal ProtocolData",
	},
	33: {
		number:   33,
		keyword:  "DCCP",
		protocol: "Datagram Congestion Control ProtocolData",
	},
	34: {
		number:   34,
		keyword:  "3PC",
		protocol: "Third Party Connect ProtocolData",
	},
	35: {
		number:   35,
		keyword:  "IDPR",
		protocol: "Inter-Domain Policy Routing ProtocolData",
	},
	36: {
		number:   36,
		keyword:  "XTP",
		protocol: "XTP",
	},
	37: {
		number:   37,
		keyword:  "DDP",
		protocol: "Datagram Delivery ProtocolData",
	},
	38: {
		number:   38,
		keyword:  "IDPR-CMTP",
		protocol: "IDPR Control Message Transport Proto",
	},
	39: {
		number:   39,
		keyword:  "TP++",
		protocol: "TP++ Transport ProtocolData",
	},
	40: {
		number:   40,
		keyword:  "IL",
		protocol: "IL Transport ProtocolData",
	},
	41: {
		number:   41,
		keyword:  "IPv6",
		protocol: "IPv6 encapsulation",
	},
	42: {
		number:   42,
		keyword:  "SDRP",
		protocol: "Source Demand Routing ProtocolData",
	},
	43: {
		number:   43,
		keyword:  "IPv6-Route",
		protocol: "Routing Header for IPv6",
	},
	44: {
		number:   44,
		keyword:  "IPv6-Frag",
		protocol: "Fragment Header for IPv6",
	},
	45: {
		number:   45,
		keyword:  "IDRP",
		protocol: "Inter-Domain Routing ProtocolData",
	},
	46: {
		number:   46,
		keyword:  "RSVP",
		protocol: "Reservation ProtocolData",
	},
	47: {
		number:   47,
		keyword:  "GRE",
		protocol: "Generic Routing Encapsulation",
	},
	48: {
		number:   48,
		keyword:  "DSR",
		protocol: "Dynamic Source Routing ProtocolData",
	},
	49: {
		number:   49,
		keyword:  "BNA",
		protocol: "BNA",
	},
	50: {
		number:   50,
		keyword:  "ESP",
		protocol: "Encap Security Payload",
	},
	51: {
		number:   51,
		keyword:  "AH",
		protocol: "Authentication Header",
	},
	52: {
		number:   52,
		keyword:  "I-NLSP",
		protocol: "Integrated Net Layer Security TUBA",
	},
	53: {
		number:   53,
		keyword:  "SWIPE (deprecated)",
		protocol: "IP with Encryption",
	},
	54: {
		number:   54,
		keyword:  "NARP",
		protocol: "NBMA Address Resolution ProtocolData",
	},
	55: {
		number:   55,
		keyword:  "Min-IPv4",
		protocol: "Minimal IPv4 Encapsulation",
	},
	56: {
		number:   56,
		keyword:  "TLSP",
		protocol: "Transport Layer Security ProtocolData using Kryptonet key management",
	},
	57: {
		number:   57,
		keyword:  "SKIP",
		protocol: "SKIP",
	},
	58: {
		number:   58,
		keyword:  "IPv6-ICMP",
		protocol: "ICMP for IPv6",
	},
	59: {
		number:   59,
		keyword:  "IPv6-NoNxt",
		protocol: "No Next Header for IPv6",
	},
	60: {
		number:   60,
		keyword:  "IPv6-Opts",
		protocol: "Destination Options for IPv6",
	},
	61: {
		number:   61,
		keyword:  "",
		protocol: "any host internal protocol",
	},
	62: {
		number:   62,
		keyword:  "CFTP",
		protocol: "CFTP",
	},
	63: {
		number:   63,
		keyword:  "",
		protocol: "any local network",
	},
	64: {
		number:   64,
		keyword:  "SAT-EXPAK",
		protocol: "SATNET and Backroom EXPAK",
	},
	65: {
		number:   65,
		keyword:  "KRYPTOLAN",
		protocol: "Kryptolan",
	},
	66: {
		number:   66,
		keyword:  "RVD",
		protocol: "MIT Remote Virtual Disk ProtocolData",
	},
	67: {
		number:   67,
		keyword:  "IPPC",
		protocol: "Internet Pluribus Packet Core",
	},
	68: {
		number:   68,
		keyword:  "",
		protocol: "any distributed file system",
	},
	69: {
		number:   69,
		keyword:  "SAT-MON",
		protocol: "SATNET Monitoring",
	},
	70: {
		number:   70,
		keyword:  "VISA",
		protocol: "VISA ProtocolData",
	},
	71: {
		number:   71,
		keyword:  "IPCV",
		protocol: "Internet Packet Core Utility",
	},
	72: {
		number:   72,
		keyword:  "CPNX",
		protocol: "Computer ProtocolData Network Executive",
	},
	73: {
		number:   73,
		keyword:  "CPHB",
		protocol: "Computer ProtocolData Heart Beat",
	},
	74: {
		number:   74,
		keyword:  "WSN",
		protocol: "Wang Span Network",
	},
	75: {
		number:   75,
		keyword:  "PVP",
		protocol: "Packet Video ProtocolData",
	},
	76: {
		number:   76,
		keyword:  "BR-SAT-MON",
		protocol: "Backroom SATNET Monitoring",
	},
	77: {
		number:   77,
		keyword:  "SUN-ND",
		protocol: "SUN ND PROTOCOL-Temporary",
	},
	78: {
		number:   78,
		keyword:  "WB-MON",
		protocol: "WIDEBAND Monitoring",
	},
	79: {
		number:   79,
		keyword:  "WB-EXPAK",
		protocol: "WIDEBAND EXPAK",
	},
	80: {
		number:   80,
		keyword:  "ISO-IP",
		protocol: "ISO Internet ProtocolData",
	},
	81: {
		number:   81,
		keyword:  "VMTP",
		protocol: "VMTP",
	},
	82: {
		number:   82,
		keyword:  "SECURE-VMTP",
		protocol: "SECURE-VMTP",
	},
	83: {
		number:   83,
		keyword:  "VINES",
		protocol: "VINES",
	},
	84: {
		number:   84,
		keyword:  "IPTM",
		protocol: "Internet ProtocolData Traffic Manager",
	},
	85: {
		number:   85,
		keyword:  "NSFNET-IGP",
		protocol: "NSFNET-IGP",
	},
	86: {
		number:   86,
		keyword:  "DGP",
		protocol: "Dissimilar Gateway ProtocolData",
	},
	87: {
		number:   87,
		keyword:  "TCF",
		protocol: "TCF",
	},
	88: {
		number:   88,
		keyword:  "EIGRP",
		protocol: "EIGRP",
	},
	89: {
		number:   89,
		keyword:  "OSPFIGP",
		protocol: "OSPFIGP",
	},
	90: {
		number:   90,
		keyword:  "Sprite-RPC",
		protocol: "Sprite RPC ProtocolData",
	},
	91: {
		number:   91,
		keyword:  "LARP",
		protocol: "Locus Address Resolution ProtocolData",
	},
	92: {
		number:   92,
		keyword:  "MTP",
		protocol: "Multicast Transport ProtocolData",
	},
	93: {
		number:   93,
		keyword:  "AX.25",
		protocol: "AX.25 Frames",
	},
	94: {
		number:   94,
		keyword:  "IPIP",
		protocol: "IP-within-IP Encapsulation ProtocolData",
	},
	95: {
		number:   95,
		keyword:  "MICP (deprecated)",
		protocol: "Mobile Internetworking Control Pro.",
	},
	96: {
		number:   96,
		keyword:  "SCC-SP",
		protocol: "Semaphore Communications Sec. Pro.",
	},
	97: {
		number:   97,
		keyword:  "ETHERIP",
		protocol: "Ethernet-within-IP Encapsulation",
	},
	98: {
		number:   98,
		keyword:  "ENCAP",
		protocol: "Encapsulation Header",
	},
	99: {
		number:   99,
		keyword:  "",
		protocol: "any private encryption scheme",
	},
	100: {
		number:   100,
		keyword:  "GMTP",
		protocol: "GMTP",
	},
	101: {
		number:   101,
		keyword:  "IFMP",
		protocol: "Ipsilon Flow Management ProtocolData",
	},
	102: {
		number:   102,
		keyword:  "PNNI",
		protocol: "PNNI over IP",
	},
	103: {
		number:   103,
		keyword:  "PIM",
		protocol: "ProtocolData Independent Multicast",
	},
	104: {
		number:   104,
		keyword:  "ARIS",
		protocol: "ARIS",
	},
	105: {
		number:   105,
		keyword:  "SCPS",
		protocol: "SCPS",
	},
	106: {
		number:   106,
		keyword:  "QNX",
		protocol: "QNX",
	},
	107: {
		number:   107,
		keyword:  "A/N",
		protocol: "Active Networks",
	},
	108: {
		number:   108,
		keyword:  "IPComp",
		protocol: "IP Payload Compression ProtocolData",
	},
	109: {
		number:   109,
		keyword:  "SNP",
		protocol: "Sitara Networks ProtocolData",
	},
	110: {
		number:   110,
		keyword:  "Compaq-Peer",
		protocol: "Compaq Peer ProtocolData",
	},
	111: {
		number:   111,
		keyword:  "IPX-in-IP",
		protocol: "IPX in IP",
	},
	112: {
		number:   112,
		keyword:  "VRRP",
		protocol: "Virtual Router Redundancy ProtocolData",
	},
	113: {
		number:   113,
		keyword:  "PGM",
		protocol: "PGM Reliable Transport ProtocolData",
	},
	114: {
		number:   114,
		keyword:  "",
		protocol: "any 0-hop protocol",
	},
	115: {
		number:   115,
		keyword:  "L2TP",
		protocol: "Layer Two Tunneling ProtocolData",
	},
	116: {
		number:   116,
		keyword:  "DDX",
		protocol: "D-II Data Exchange (DDX)",
	},
	117: {
		number:   117,
		keyword:  "IATP",
		protocol: "Interactive Agent Transfer ProtocolData",
	},
	118: {
		number:   118,
		keyword:  "STP",
		protocol: "Schedule Transfer ProtocolData",
	},
	119: {
		number:   119,
		keyword:  "SRP",
		protocol: "SpectraLink Radio ProtocolData",
	},
	120: {
		number:   120,
		keyword:  "UTI",
		protocol: "UTI",
	},
	121: {
		number:   121,
		keyword:  "SMP",
		protocol: "Simple Message ProtocolData",
	},
	122: {
		number:   122,
		keyword:  "SM (deprecated)",
		protocol: "Simple Multicast ProtocolData",
	},
	123: {
		number:   123,
		keyword:  "PTP",
		protocol: "Performance Transparency ProtocolData",
	},
	124: {
		number:   124,
		keyword:  "ISIS over IPv4",
		protocol: "",
	},
	125: {
		number:   125,
		keyword:  "FIRE",
		protocol: "",
	},
	126: {
		number:   126,
		keyword:  "CRTP",
		protocol: "Combat Radio Transport ProtocolData",
	},
	127: {
		number:   127,
		keyword:  "CRUDP",
		protocol: "Combat Radio User Datagram",
	},
	128: {
		number:   128,
		keyword:  "SSCOPMCE",
		protocol: "",
	},
	129: {
		number:   129,
		keyword:  "IPLT",
		protocol: "",
	},
	130: {
		number:   130,
		keyword:  "SPS",
		protocol: "Secure Packet Shield",
	},
	131: {
		number:   131,
		keyword:  "PIPE",
		protocol: "Private IP Encapsulation within IP",
	},
	132: {
		number:   132,
		keyword:  "SCTP",
		protocol: "Stream Control Transmission ProtocolData",
	},
	133: {
		number:   133,
		keyword:  "FC",
		protocol: "Fibre Channel",
	},
	134: {
		number:   134,
		keyword:  "RSVP-E2E-IGNORE",
		protocol: "",
	},
	135: {
		number:   135,
		keyword:  "Mobility Header",
		protocol: "",
	},
	136: {
		number:   136,
		keyword:  "UDPLite",
		protocol: "",
	},
	137: {
		number:   137,
		keyword:  "MPLS-in-IP",
		protocol: "",
	},
	138: {
		number:   138,
		keyword:  "manet",
		protocol: "MANET Protocols",
	},
	139: {
		number:   139,
		keyword:  "HIP",
		protocol: "Host Identity ProtocolData",
	},
	140: {
		number:   140,
		keyword:  "Shim6",
		protocol: "Shim6 ProtocolData",
	},
	141: {
		number:   141,
		keyword:  "WESP",
		protocol: "Wrapped Encapsulating Security Payload",
	},
	142: {
		number:   142,
		keyword:  "ROHC",
		protocol: "Robust Header Compression",
	},
	143: {
		number:   143,
		keyword:  "Ethernet",
		protocol: "Ethernet",
	},
	144: {
		number:   144,
		keyword:  "AGGFRAG",
		protocol: "AGGFRAG encapsulation payload for ESP",
	},
	145: {
		number:   145,
		keyword:  "NSH",
		protocol: "Network Service Header",
	},
}
