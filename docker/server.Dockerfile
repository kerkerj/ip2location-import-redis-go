FROM alpine:3.12
WORKDIR /var/www/go
ADD ./ip2location-server-linux .
ADD ./IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ISP-DOMAIN-MOBILE-USAGETYPE.SAMPLE.BIN .
CMD '/var/www/go/ip2location-server-linux'