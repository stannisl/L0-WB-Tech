package main

import (
	"L0-wb/database"
	"L0-wb/kafka"
	"L0-wb/models"
	"log"
	"net/http"

	"L0-wb/cache"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Connecting to Database")
	database.Connect()

	log.Println("Starting Kafka Consumer")
	go kafka.StartKafkaConsumer("localhost:9092", "orders", "go-consumer-group")

	go StartServer()

	select {}
}

func StartServer() {
	r := gin.Default()

	r.GET("/order/:orderUid", func(c *gin.Context) {
		orderUid := c.Param("orderUid")
		var order models.Order

		order, ok := cache.CacheMap[orderUid]
		if !ok {
			result := database.GetDB().
				Preload("Items").
				Preload("Delivery").
				Preload("Payment").
				First(&order, "order_uid = ?", orderUid)

			if result.Error != nil {
				log.Println("Error while fetching order from DB: ", result.Error)
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
			cache.CacheMap[orderUid] = order
		}

		c.JSON(http.StatusOK, order)
	})

	r.Run(":8080")
}
