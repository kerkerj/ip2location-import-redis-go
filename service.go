package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/getsocial-rnd/ip2location-go"
	"github.com/go-redis/redis/v8"
)

type GeoIPService struct {
	redisClient   *redis.Client
	ip2locationDB *ip2location.DB
	redisKey string
}

func NewGeoIPService(binaryDBPath, redisDSN, redisKey string) *GeoIPService {
	if binaryDBPath == "" {
		binaryDBPath = TestDBPath
	}
	if redisDSN == "" {
		redisDSN = DefaultRedisDSN
	}
	if redisKey == "" {
		redisKey = DefaultRedisKey
	}

	db, openDBError := ip2location.Open(binaryDBPath)
	if openDBError != nil {
		log.Fatal(openDBError)
	}

	return &GeoIPService{
		redisClient: redis.NewClient(&redis.Options{
			Network:    "",
			Addr:       redisDSN,
			Password:   "",
			DB:         0,
			MaxRetries: 5,
			PoolSize:   10,
		}),
		ip2locationDB: db,
		redisKey: DefaultRedisKey,
	}
}

func (g *GeoIPService) SetRedisKey(key string) *GeoIPService {
	g.redisKey = key
	return g
}

func (g *GeoIPService) GetIPInfoFromBin(ipAddress string) *ip2location.Record {
	ipBigInt := IPtoBigInt(net.ParseIP(ipAddress))
	r, err := g.redisClient.ZRevRangeByScore(context.Background(), g.redisKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    ipBigInt.String(),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Fatal("GetIPInfoFromBin: ", err)
	}

	if len(r) == 0 {
		return nil
	}

	return parseToRecord(r[0])
}

func (g *GeoIPService) GetIPInfoFromRedis(ipAddress string) *ip2location.Record {
	ipBigInt := IPtoBigInt(net.ParseIP(ipAddress))
	r, err := g.redisClient.ZRevRangeByScore(context.Background(), g.redisKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    ipBigInt.String(),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Fatal("GetIPInfoFromRedis: ", err)
	}

	if len(r) == 0 {
		return nil
	}

	return parseToRecord(r[0])
}

func (g *GeoIPService) GetEncodedIPInfoFromRedis(ipAddress string) *ip2location.Record {
	ipBigInt := IPtoBigInt(net.ParseIP(ipAddress))
	r, err := g.redisClient.ZRevRangeByScore(context.Background(), g.redisKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    ipBigInt.String(),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Fatal("GetEncodedIPInfoFromRedis: ", err)
	}

	if len(r) == 0 {
		return nil
	}

	decoded := GobDecode(r[0])
	return &decoded.R
}

func parseToRecord(str string) *ip2location.Record {
	// IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ISP-DOMAIN-MOBILE-USAGETYPE
	// 281473506541568|TW|Taiwan|Taipei|Taipei|25.047760|121.531850|Chunghwa Telecom Co. Ltd.|cht.com.tw|466|11/92|Chunghwa LDM|ISP/MOB
	strArray := strings.Split(str, "|")
	lat, _ := strconv.ParseFloat(strArray[5], 32)
	lon, _ := strconv.ParseFloat(strArray[6], 32)
	return &ip2location.Record{
		CountryShort:       strArray[1],
		CountryLong:        strArray[2],
		Region:             strArray[3],
		City:               strArray[4],
		Isp:                strArray[7],
		Latitude:           float32(lat),
		Longitude:          float32(lon),
		Domain:             strArray[8],
		Zipcode:            "",
		TimeZone:           "",
		NetSpeed:           "",
		IddCode:            "",
		Areacode:           "",
		WeatherStationCode: "",
		WeatherStationName: "",
		Mcc:                strArray[9],
		Mnc:                strArray[10],
		MobileBrand:        strArray[11],
		Elevation:          0,
		UsageType:          strArray[12],
	}
}

