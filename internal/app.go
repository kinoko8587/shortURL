package internal

import (
	"time"

	"example.com/greetings/internal/controller"
	"example.com/greetings/internal/middleware"
	"github.com/gin-gonic/gin"
)

func App() {
	router := gin.Default()

	router.Use(middleware.RateLimitMiddlewareByIPAndUser(
		10,          // bucket capacity
		5,           // refill rate
		time.Second, // refill every second
	))

	controller.Controller(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.Run(":8080")
}
