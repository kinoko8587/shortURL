package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v4"
)

var (
	urlStore = sync.Map{}
)

func Store(shortURL string, longURL string) {
	urlStore.Store(shortURL, longURL)
}

func Get(shortURL string) (string, bool) {
	value, ok := urlStore.Load(shortURL)
	return value.(string), ok
}

func connectDB() {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	fmt.Println("Connected to PostgreSQL successfully!")
}
