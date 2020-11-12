package util

import (
	"crypto/sha1"
	"encoding/hex"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adriansr/flowhash"
)

// GetEnv Handel os.LookupEnv
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func Hash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func RemoveAfter(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func RemoveBefor(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func ParseBool(s string) bool {
	b, _ := strconv.ParseBool(s)

	return b
}

func CommunityIDHash(srcIP, dstIP, srcPort, dstPort, protocol string) string {
	const (
		icmpProtocol     uint8 = 1
		igmpProtocol     uint8 = 2
		tcpProtocol      uint8 = 6
		udpProtocol      uint8 = 17
		greProtocol      uint8 = 47
		icmpIPv6Protocol uint8 = 58
		eigrpProtocol    uint8 = 88
		ospfProtocol     uint8 = 89
		pimProtocol      uint8 = 103
		sctpProtocol     uint8 = 132
	)

	var transports = map[string]uint8{
		"icmp":      icmpProtocol,
		"igmp":      igmpProtocol,
		"tcp":       tcpProtocol,
		"udp":       udpProtocol,
		"gre":       greProtocol,
		"ipv6-icmp": icmpIPv6Protocol,
		"icmpv6":    icmpIPv6Protocol,
		"eigrp":     eigrpProtocol,
		"ospf":      ospfProtocol,
		"pim":       pimProtocol,
		"sctp":      sctpProtocol,
	}
	sPort, _ := strconv.Atoi(srcPort)
	dPort, _ := strconv.Atoi(dstPort)
	flow := flowhash.Flow{
		SourceIP:        net.ParseIP(srcIP),
		DestinationIP:   net.ParseIP(dstIP),
		SourcePort:      uint16(sPort),
		DestinationPort: uint16(dPort),
		Protocol:        transports[protocol],
	}
	return flowhash.CommunityID.Hash(flow)
}
