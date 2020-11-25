package main

import (
	"encoding/csv"
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestIPtoBigInt(t *testing.T) {
	tests := []struct {
		name string
		ip string
		want string
	}{
		{
			name: "1.1.1.1",
			ip: CloudFlareDNSIPv4,
			want: "281470698586369",
		},
		{
			name: "2606:4700:4700::1111 (IPv6 of 1.1.1.1)",
			ip: CloudFlareDNSIPv6,
			want: "50543257694033307102031451402929180945",
		},
		{
			name: "168.95.1.1",
			ip: HiNetDNSIPv4,
			want: "281473506541825",
		},
		{
			name: "2001:b000:168::1 (IPv6 of 168.95.1.1)",
			ip: HiNetDNSIPv6,
			want: "42544057866501298749606237644420808705",
		},
		{
			name: "8.8.8.8",
			ip: GoogleDNSIpv4,
			want: "281470816487432",
		},
		{
			name: "2001:4860:4860::8888 (IPv6 of 8.8.8.8)",
			ip: GoogleDNSIpv6,
			want: "42541956123769884636017138956568135816",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPtoBigInt(net.ParseIP(tt.ip)); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("IPtoBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessLocationLine(t *testing.T) {
	// Arrange
	expected := "281470698520576|US|United States of America|California|Los Angeles|34.052230|-118.243680|APNIC and CloudFlare DNS Resolver Project|cloudflare.com|-|-|-|CDN"
	rawData := `"281470698520576","281470698520831","US","United States of America","California","Los Angeles","34.052230","-118.243680","APNIC and CloudFlare DNS Resolver Project","cloudflare.com","-","-","-","CDN"`

	// Act
	csvData, _ := csv.NewReader(strings.NewReader(rawData)).ReadAll()
	processed := ProcessLocationLine(csvData[0])

	t.Log(csvData)
	t.Log(expected)
	t.Log(processed)

	// Assert
	if processed != expected {
		t.Error("not expected")
	}
}
