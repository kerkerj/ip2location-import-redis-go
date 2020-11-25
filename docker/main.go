package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsocial-rnd/ip2location-go"
)

const (
	dbPath = "./IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ISP-DOMAIN-MOBILE-USAGETYPE.SAMPLE.BIN"
)

// curl http://localhost:8080?ip=127.0.0.1
func main() {
	envDBPath := os.Getenv("DB")
	if envDBPath == "" {
		envDBPath = dbPath
	}

	// open db
	db, err := ip2location.Open(envDBPath)
	if err != nil {
		log.Fatal(err)
	}

	handleFunc := func(w http.ResponseWriter, req *http.Request) {
		ips, ok := req.URL.Query()["ip"]

		if !ok || len(ips[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Url Param 'key' is missing"))
			return
		}

		// get data from db
		record, err := db.GetAll(ips[0])
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("IP: %s", record.String())))
		return
	}

	http.HandleFunc("/", handleFunc)
	log.Print("starting server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
