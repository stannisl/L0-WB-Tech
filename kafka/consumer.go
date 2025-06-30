package kafka

import (
	"L0-wb/database"
	"L0-wb/models"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func StartKafkaConsumer(brokerAddress, topic, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{

		Brokers:           []string{brokerAddress},
		Topic:             topic,
		GroupID:           groupID,
		CommitInterval:    0,
		MaxWait:           2 * time.Second,
		StartOffset:       kafka.FirstOffset,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
		RetentionTime:     24 * time.Hour,
	})

	log.Println("Kafka consumer started successfully")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var order models.Order

		if err = json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			if err := r.CommitMessages(context.Background(), msg); err != nil {
				log.Printf("failed commit: %v", err)
			}
			continue
		}

		db := database.GetDB()

		var existingOrder models.Order
		if err := db.Where("order_uid = ?", order.OrderUid).First(&existingOrder).Error; err == nil {
			log.Printf("error adding record from Kafka with offset = %v. Order with uid = %s already exists", msg.Offset, order.OrderUid)
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("db error: %v", err)
			continue
		}

		tx := db.Begin()
		if err = tx.Create(&order).Error; err != nil {
			log.Printf("error creating order: %v", err)
			tx.Rollback()
			continue
		}
		tx.Commit()

		if err := r.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("failed commit: %v", err)
		}

		log.Printf("Successfully processed and committed message (offset: %d, order_uid: %s)", msg.Offset, order.OrderUid)
	}
}
