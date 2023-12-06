package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5434
	user     = "user"
	password = "123"
	dbname   = "test"
)

// var query = `SELECT country, star_sign, count(*) FROM users
// WHERE balance BETWEEN 301 AND 700 GROUP BY country, star_sign;`

// var query = `SELECT country, star_sign, count(*) FROM users
// WHERE balance_range = 2 GROUP BY country, star_sign;`

var query = `SELECT country, star_sign, avg(balance) FROM users GROUP BY country, star_sign;`

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	loops := 10
	totalDuration := time.Duration(0)
	for i < loops {
		start := time.Now()
		r, err := db.Query(query)

		if err != nil {
			log.Fatal(err)
		}

		records := 0
		for r.Next() {
			var country string
			var star_sign string
			var avg float64
			err = r.Scan(&country, &star_sign, &avg)
			if err != nil {
				log.Fatal(err)
			}
			records++
		}
		duration := time.Since(start)
		totalDuration += duration
		fmt.Printf("-- [Round %v] %v records - %v\n", i, records, duration)
		i++
	}

	fmt.Printf("==>> [Average] %v", totalDuration/time.Duration(loops))

}
