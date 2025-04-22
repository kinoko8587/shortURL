package usecases

import (
	"net/http"

	"example.com/greetings/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

func RedirectShortURL(c *gin.Context) {
	shortURL := c.Param("short_url")
	longURL, ok := infrastructure.Get(shortURL)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, longURL)
}
