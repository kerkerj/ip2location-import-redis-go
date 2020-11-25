# Import ip2location DB data to Redis

Base on the article [Importing IP2Location data into Redis and querying with PHP](https://blog.ip2location.com/knowledge-base/importing-ip2location-data-into-redis-and-querying-with-php/), the original script is written in perl, in this repo we rewrite the script in Golang, and use [ip2location DB23 (IPv6)](https://www.ip2location.com/database/db23-ip-country-region-city-latitude-longitude-isp-domain-mobile-usagetype) for examples.

Sample db files can be found here: [https://www.ip2location.com/development-libraries](https://www.ip2location.com/development-libraries)

Test code is in `main.go`, it demonstrates how to import data and read it.