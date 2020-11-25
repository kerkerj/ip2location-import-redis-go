package main

import "fmt"

const (
	DefaultRedisKey = "DB23"
	DefaultRedisDSN = "localhost:6379"
	TestCSVPath = "./testdata/IPV6.DB23.SAMPLE.csv"
	TestDBPath = "./testdata/IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ISP-DOMAIN-MOBILE-USAGETYPE.SAMPLE.BIN"

	CloudFlareDNSIPv4 = "1.1.1.1" // only this one is in test data
	CloudFlareDNSIPv6 = "2606:4700:4700::1111"
	GoogleDNSIpv4     = "8.8.8.8"
	GoogleDNSIpv6     = "2001:4860:4860::8888"
	HiNetDNSIPv4      = "168.95.1.1"
	HiNetDNSIPv6      = "2001:b000:168::1"
)

// For testing
func main() {
	// Import csv to redis, you can choose to import encoded data or not
	importService := NewImportService(DefaultRedisKey, TestCSVPath, DefaultRedisDSN)
	importService.SetRedisKey("DB234").ImportData()
	importService.SetRedisKey("DB23Encoded").ImportEncodedData()

	// Test query
	geoIPService := NewGeoIPService(TestDBPath, DefaultRedisDSN, DefaultRedisKey)
	result1 := geoIPService.GetIPInfoFromBin(CloudFlareDNSIPv4)
	result2 := geoIPService.SetRedisKey("DB234").GetIPInfoFromRedis(CloudFlareDNSIPv4)
	result3 := geoIPService.SetRedisKey("DB23Encoded").GetEncodedIPInfoFromRedis(CloudFlareDNSIPv4)

	fmt.Printf("%+v\n", result1)
	fmt.Println("---")
	fmt.Printf("%+v\n", result2)
	fmt.Println("---")
	fmt.Printf("%+v\n", result3)
}