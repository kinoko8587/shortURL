package usecases

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"

	"example.com/greetings/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

func CreateShortURL(c *gin.Context) {
	longURL := c.PostForm("long_url")
	if longURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "long_url is required"})
		return
	}
	shortURL := generateShortURL(longURL)
	infrastructure.Store(shortURL, longURL)

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortURL,
		"long_url":  longURL,
	})
}

func generateShortURL(longURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(longURL))
	shortURL := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:6]
	return shortURL
}
