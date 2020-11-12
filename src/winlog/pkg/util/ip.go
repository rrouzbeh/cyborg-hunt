package util

import (
	"strconv"
	"strings"
)

func IPProcess(s string) (string, string, string, string) {
	rfc := ""
	isPublic := ""
	ipType := ""
	ipLen := len(s)
	ipVersion := "4"
	if !strings.Contains(s, ":") && ipLen <= 15 && ipLen >= 7 {
		if strings.HasPrefix(s, "10.") {
			rfc = "RFC_1918"
			isPublic = "false"
			ipType = "private"
		} else if strings.HasPrefix(s, "192.168.") {
			rfc = "RFC_1918"
			isPublic = "false"
			ipType = "private"
		} else if strings.HasPrefix(s, "169.254.") {
			rfc = "RFC_3927"
			isPublic = "false"
			ipType = "local"
		} else if strings.HasPrefix(s, "127.") {
			rfc = "RFC_1122-3.2.1.3"
			isPublic = "false"
			ipType = "loopback"
		} else if strings.HasPrefix(s, "0.") {
			rfc = "RFC_1700"
			isPublic = "false"
			ipType = "reserved_as_a_source_address_only"
		} else if strings.HasPrefix(s, "192.88.99.") {
			rfc = "RFC_3068"
			isPublic = "false"
			ipType = "6to4"
		} else if strings.HasPrefix(s, "192.31.196.") {
			rfc = "RFC_3068"
			isPublic = "false"
			ipType = "as112-v4"
		} else if strings.HasPrefix(s, "192.52.193.") {
			rfc = "RFC_7450"
			isPublic = "false"
			ipType = "amt"
		} else if strings.HasPrefix(s, "172.") {
			o := strings.Split(s, ".")
			n, _ := strconv.Atoi(o[1])
			if n >= 16 && n <= 31 {
				rfc = "RFC_1918"
				isPublic = "false"
				ipType = "private"
			} else {
				rfc = "RFC_1366"
				isPublic = "true"
				ipType = "public"
			}
		} else if strings.HasPrefix(s, "100.") {
			o := strings.Split(s, ".")
			n, _ := strconv.Atoi(o[1])
			if n >= 64 && n <= 127 {
				rfc = "RFC_1918"
				isPublic = "false"
				ipType = "private"
			} else {
				rfc = "RFC_1366"
				isPublic = "true"
				ipType = "public"
			}
		} else if strings.HasPrefix(s, "2") {
			o := strings.Split(s, ".")
			n, _ := strconv.Atoi(o[0])
			if s == "255.255.255.255" {
				rfc = "RFC_8190"
				isPublic = "false"
				ipType = "broadcast"
			} else if n >= 224 && n <= 255 {
				rfc = "RFC_1112"
				isPublic = "false"
				ipType = "multicast"
			} else {
				rfc = "RFC_1366"
				isPublic = "true"
				ipType = "public"
			}
		} else {
			rfc = "RFC_1366"
			isPublic = "true"
			ipType = "public"
		}
	} else {
		rfc = "n/a"
		isPublic = "false"
		ipType = "n/a"
		ipVersion = "6"
	}
	return rfc, isPublic, ipType, ipVersion
}
