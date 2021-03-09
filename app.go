package main

import (
	"golang-azure-eventhub/adapter"
	"golang-azure-eventhub/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	eventHubConnection := adapter.NewConnection()
	log.Println("Conectado", eventHubConnection)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", controllers.HealthControllerHandler())
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"msg": "Not found"})
	})

	router.Run(":8080")
}
