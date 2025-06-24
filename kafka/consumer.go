package kafka

import (
	"L0-wb/database"
	"L0-wb/models"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(brokerAddress, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10 кб (красное белое)
		MaxBytes:       10e6, // 10 мб
		CommitInterval: 0,
	})

	log.Println("Kafka consumer started successfully")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ error reading message: %v", err)
			continue
		}

		order := models.Order{}
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("❌ error unmarshaling message: %v", err)
			continue
		}

		db := database.GetDB()
		if err := db.Create(&order); err != nil {
			log.Printf("❌ error creating order: %v", err)
			continue
		}
	}
}
