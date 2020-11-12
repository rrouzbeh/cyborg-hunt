package util

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/common/log"
)

func Geocountry(s string) string {
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(s)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	return record.Country.IsoCode
}
func GeoASN(s string) (uint32, string) {
	db, err := geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(s)
	record, err := db.ASN(ip)
	if err != nil {
		log.Fatal(err)
	}
	return uint32(record.AutonomousSystemNumber), record.AutonomousSystemOrganization
}
