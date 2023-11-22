package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
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

var (
	numOfTickets = 100
	numOfBookers = 110
	simulateTime = 50 * time.Millisecond
)

func bookTicketWithoutSkipLocked(db *sql.DB, eventID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var ticketID int
	err = tx.QueryRow("SELECT ticket_id FROM tickets WHERE event_id = $1 AND status = 'available' LIMIT 1 FOR UPDATE", eventID).Scan(&ticketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	time.Sleep(simulateTime)

	_, err = tx.Exec("UPDATE tickets SET status = 'sold' WHERE ticket_id = $1", ticketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func bookTicketWithSkipLocked(db *sql.DB, eventID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var ticketID int
	err = tx.QueryRow("SELECT ticket_id FROM tickets WHERE event_id = $1 AND status = 'available' LIMIT 1 FOR UPDATE SKIP LOCKED", eventID).Scan(&ticketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	time.Sleep(simulateTime)

	_, err = tx.Exec("UPDATE tickets SET status = 'sold' WHERE ticket_id = $1", ticketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

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

	initTables(db)

	seedTables(db)

	fmt.Println("Booking tickets with SKIP LOCKED")
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < numOfBookers; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("Booking ticket %d\n", i)
			defer wg.Done()
			err := bookTicketWithoutSkipLocked(db, 1)
			// err := bookTicketWithSkipLocked(db, 2)
			if err != nil {
				log.Printf("Failed to book ticket: %v\n", err)
			} else {
				log.Printf("Successfully booked ticket %d\n", i)
			}
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	log.Printf("Total time for booking without SKIP LOCKED: %s", duration)
}

func initTables(db *sql.DB) {
	var err error

	_, err = db.Exec(`
		DROP TABLE IF EXISTS events CASCADE;
		CREATE TABLE events (
			event_id SERIAL PRIMARY KEY,
			name VARCHAR(255)
		);
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		DROP TABLE IF EXISTS tickets CASCADE;
		CREATE TABLE tickets (
			ticket_id SERIAL PRIMARY KEY,
			event_id INT NOT NULL,
			status VARCHAR(50) NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func seedTables(db *sql.DB) {
	var err error

	_, err = db.Exec(`
		INSERT INTO events (event_id, name)
		VALUES (1, 'Event 1'),
			   (2, 'Event 2');
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf(`
		INSERT INTO tickets (event_id, status)
		SELECT 1, 'available'
		FROM generate_series(1, %d);
	`, numOfTickets))

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf(`
		INSERT INTO tickets (event_id, status)
		SELECT 2, 'available'
		FROM generate_series(1, %d);
	`, numOfTickets))

	if err != nil {
		log.Fatal(err)
	}
}
