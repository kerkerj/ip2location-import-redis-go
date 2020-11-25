package main

import (
	"testing"

	"github.com/getsocial-rnd/ip2location-go"
)

func BenchmarkReadFromBinary(b *testing.B) {
	s := NewGeoIPService(TestDBPath, DefaultRedisDSN, DefaultRedisKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ip2locationDB.GetAll(CloudFlareDNSIPv4)
	}
}

func BenchmarkReadFromRedis(b *testing.B) {
	s := NewGeoIPService(TestDBPath, DefaultRedisDSN, DefaultRedisKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetIPInfoFromRedis(CloudFlareDNSIPv4)
	}
}

func BenchmarkReadEncodedFromRedis(b *testing.B) {
	s := NewGeoIPService(TestDBPath, DefaultRedisDSN, DefaultRedisKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetEncodedIPInfoFromRedis(CloudFlareDNSIPv4)
	}
}

func BenchmarkStringToIP2LocationStruct(b *testing.B) {
	str := "281470698521600|AU|Australia|Victoria|Melbourne|-37.814000|144.963320|WirefreeBroadband Pty Ltd|wirefreebroadband.com.au|-|-|-|ISP"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseToRecord(str)
	}
}

func BenchmarkGobDecodeToCustomStruct(b *testing.B) {
	data := M{
		S: 1,
		R: ip2location.Record{
			CountryShort:       "AU",
			CountryLong:        "Australia",
			Region:             "Victoria",
			City:               "Melbourne",
			Isp:                "WirefreeBroadband Pty Ltd",
			Latitude:           -37.814000,
			Longitude:          144.963320,
			Domain:             "wirefreebroadband.com.au",
			Zipcode:            "",
			TimeZone:           "",
			NetSpeed:           "",
			IddCode:            "",
			Areacode:           "",
			WeatherStationCode: "",
			WeatherStationName: "",
			Mcc:                "-",
			Mnc:                "-",
			MobileBrand:        "",
			Elevation:          0,
			UsageType:          "ISP",
		},
	}

	encodedStr := GobEncode(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GobDecode(encodedStr)
	}
}

