package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	_ "github.com/lib/pq"
)

// Flags for database connection, concurrency, and iterations
var (
	connString  = flag.String("conn", "postgres:///postgres?sslmode=disable", "Postgres connection string")
	concurrency = flag.Int("concurrency", 10, "Number of concurrent requests")
	iterations  = flag.Int("iterations", 5, "Number of updates per goroutine")
)

func main() {
	flag.Parse()

	// Connect to the database
	db, err := sql.Open("postgres", *connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Recreate the table and insert a row
	err = setupTable(db)
	if err != nil {
		log.Fatalf("Failed to setup table: %v", err)
	}

	// Initialize histogram to measure latencies in microseconds
	hist := hdrhistogram.New(1, 10000000, 3) // 1us to 10s with 3 significant figures

	// Run the update concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			updateRows(db, *iterations, hist, &mu)
		}(i)
	}
	wg.Wait()

	// Output histogram data
	fmt.Printf("\nLatency Distribution (microseconds):\n")
	fmt.Printf("Min: %v µs\n", hist.Min())
	fmt.Printf("Max: %v µs\n", hist.Max())
	fmt.Printf("Mean: %v µs\n", hist.Mean())
	fmt.Printf("P50: %v µs\n", hist.ValueAtQuantile(50))
	fmt.Printf("P90: %v µs\n", hist.ValueAtQuantile(90))
	fmt.Printf("P99: %v µs\n", hist.ValueAtQuantile(99))
}

// setupTable drops and recreates the 'foo' table, then inserts a single row
func setupTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS foo`)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE foo (id SERIAL PRIMARY KEY, updated_at TIMESTAMP DEFAULT now())`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	_, err = db.Exec(`INSERT INTO foo (updated_at) VALUES (now())`)
	if err != nil {
		return fmt.Errorf("error inserting initial row: %w", err)
	}

	log.Println("Table foo recreated and initialized with one row.")
	return nil
}

// updateRows performs the updates and tracks latency for each one
func updateRows(db *sql.DB, iterations int, hist *hdrhistogram.Histogram, mu *sync.Mutex) {
	stmt, err := db.Prepare(`UPDATE foo SET updated_at = now() WHERE id = 1`)
	ctx := context.Background()

	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, err := stmt.ExecContext(ctx)
		if err != nil {
			log.Printf("Error updating row: %v", err)
		}
		latency := time.Since(start).Microseconds()

		// Lock histogram to safely record latency
		mu.Lock()
		hist.RecordValue(latency)
		mu.Unlock()
	}
}
