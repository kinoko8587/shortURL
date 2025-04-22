package internal

import (
	"example.com/greetings/internal/controller"
	"github.com/gin-gonic/gin"
)

func App() {
	router := gin.Default()
	controller.Controller(router)
	router.Run(":8080")
}
