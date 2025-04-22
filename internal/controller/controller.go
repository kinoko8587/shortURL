package controller

import (
	"example.com/greetings/internal/usecases"
	"github.com/gin-gonic/gin"
)

func Controller(router *gin.Engine) {
	router.POST("/shorten", usecases.CreateShortURL)
	router.GET("/:short_url", usecases.RedirectShortURL)
}
