package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/getsocial-rnd/ip2location-go"
	"github.com/go-redis/redis/v8"
)

type ImportService struct {
	redisClient    *redis.Client
	csvPath        string
	importRedisKey string
}

func NewImportService(redisKey, csvPath, redisDSN string) *ImportService {
	if redisKey == "" {
		redisKey = DefaultRedisKey
	}
	if csvPath == "" {
		csvPath = TestCSVPath
	}
	if redisDSN == "" {
		redisDSN = DefaultRedisDSN
	}

	return &ImportService{
		redisClient: redis.NewClient(&redis.Options{
			Addr:       redisDSN,
			Password:   "",
			DB:         0,
			MaxRetries: 5,
			PoolSize:   10,
		}),
		csvPath:        csvPath,
		importRedisKey: redisKey,
	}
}

func (s *ImportService) SetRedisKey(key string) *ImportService {
	s.importRedisKey = key
	return s
}

func (s *ImportService) ImportEncodedData() {
	s.openCSVAndSaveToRedis(true)
}

func (s *ImportService) ImportData() {
	s.openCSVAndSaveToRedis(false)
}

func (s *ImportService) writeDataToRedis(data []*redis.Z) error {
	return s.redisClient.ZAdd(context.Background(), s.importRedisKey, data...).Err()
}

func (s *ImportService) openCSVAndSaveToRedis(saveEncoded bool) {
	// if saveEncoded = true, store GobEncoded data to redis
	generateFunc := generateCandidate
	if saveEncoded {
		generateFunc = generateEncodedCandidate
	}

	// open csv
	csvFile, err := os.Open(s.csvPath)
	if err != nil {
		log.Fatal("open csv", err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	// read data from CSV and write to redis
	r := csv.NewReader(csvFile)
	counter := 0
	candidatesData := make([]*redis.Z, 0)
	for {
		counter++

		record, err := r.Read()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			log.Fatal("csv read", err)
		}

		// validate the record
		if len(record) != 14 {
			log.Fatal(errors.New(fmt.Sprintf("line is invalid: %v", record)))
		}

		candidatesData = append(candidatesData, generateFunc(record))

		// add to redis every 5000 records
		if counter%5000 == 0 {
			if err := s.writeDataToRedis(candidatesData); err != nil {
				log.Fatal("openCSVAndSaveToRedis", err)
			}
			fmt.Printf("Running line %d\n", counter)

			// reset
			candidatesData = make([]*redis.Z, 0)
		}
	}

	// insert the rest
	if len(candidatesData) != 0 {
		if err := s.writeDataToRedis(candidatesData); err != nil {
			log.Fatal("openCSVAndSaveToRedis", err)
		}
	}
}

func generateEncodedCandidate(record []string) *redis.Z {
	score, _ := strconv.ParseFloat(record[0], 64)
	parsedRecord := StrArrayToIP2LocationRecord(score, record)
	return 	&redis.Z{
		Score: score,
		Member: GobEncode(parsedRecord),
	}
}

func generateCandidate(record []string) *redis.Z {
	score, _ := strconv.ParseFloat(record[0], 64)
	return 	&redis.Z{
		Score: score,
		Member: ProcessLocationLine(record),
	}
}

// M is a wrapped struct of ip2location, it stores additional score in order to prevent collision
type M struct {
	S float64 // avoid collision in redis
	R ip2location.Record
}

func GobEncode(data M) string {
	var network bytes.Buffer
	gob.Register(M{})
	enc := gob.NewEncoder(&network)
	err := enc.Encode(data)
	if err != nil {
		log.Fatal("GobEncode: ", err)
	}
	return network.String()
}

func GobDecode(dataStr string) M {
	var target M
	encodedReceived := []byte(dataStr)
	encodedReceivedIOReader := bytes.NewBuffer(encodedReceived)
	if err := gob.NewDecoder(encodedReceivedIOReader).Decode(&target); err != nil {
		log.Fatal("GobDecode: ", err, "; data: ", dataStr)
	}
	return target
}

func IPtoBigInt(IPAddress net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPAddress.To16())
	return IPv6Int
}

// raw data example:
// 	"281470698521600","281470698522623","AU","Australia","Victoria","Melbourne","-37.814000","144.963320","WirefreeBroadband Pty Ltd","wirefreebroadband.com.au","-","-","-","ISP"
// redis:
//	score: 281470698521600
//	member: "281470698521600","AU","Australia","Victoria","Melbourne","-37.814000","144.963320","WirefreeBroadband Pty Ltd","wirefreebroadband.com.au","-","-","-","ISP"
func StrArrayToIP2LocationRecord(start float64, strArray []string) M {
	lat, _ := strconv.ParseFloat(strArray[6], 32)
	lon, _ := strconv.ParseFloat(strArray[7], 32)
	return M{
		S: start,
		R: ip2location.Record{
			CountryShort:       strArray[2],
			CountryLong:        strArray[3],
			Region:             strArray[4],
			City:               strArray[5],
			Isp:                strArray[8],
			Latitude:           float32(lat),
			Longitude:          float32(lon),
			Domain:             strArray[9],
			Zipcode:            "",
			TimeZone:           "",
			NetSpeed:           "",
			IddCode:            "",
			Areacode:           "",
			WeatherStationCode: "",
			WeatherStationName: "",
			Mcc:                strArray[10],
			Mnc:                strArray[11],
			MobileBrand:        strArray[12],
			Elevation:          0,
			UsageType:          strArray[13],
		},
	}
}

// transfer this:
// 	`"281470698520576","281470698520831","US","United States of America","California","Los Angeles","34.052230","-118.243680","APNIC and CloudFlare DNS Resolver Project","cloudflare.com","-","-","-","CDN"`
// to:
// 	"281470698520576|US|United States of America|California|Los Angeles|34.052230|-118.243680|APNIC and CloudFlare DNS Resolver Project|cloudflare.com|-|-|-|CDN"
func ProcessLocationLine(line []string) string {
	// remove second element
	target := append(line[:1], line[2:]...)
	return strings.Join(target, "|")
}
