package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Docker!")
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

var (
	urlStore = sync.Map{}
)

func generateShortURL(longURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(longURL))
	shortURL := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:6]
	return shortURL
}

func createShortURL(c *gin.Context) {
	longURL := c.PostForm("long_url")
	if longURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "long_url is required"})
		return
	}
	shortURL := generateShortURL(longURL)
	urlStore.Store(shortURL, longURL)

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortURL,
		"long_url":  longURL,
	})
}

func redirectShortURL(c *gin.Context) {
	shortURL := c.Param("short_url")
	value, ok := urlStore.Load(shortURL)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	longURL := value.(string)
	c.Redirect(http.StatusMovedPermanently, longURL)
}

func main() {
	router := gin.Default()

	router.POST("/shorten", createShortURL)
	router.GET("/:short_url", redirectShortURL)

	router.Run(":8080")
}
