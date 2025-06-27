package service

import (
	"L0-wb/cache"
	"L0-wb/database"
	"L0-wb/models"
	"log"
)

func SaveToCache(orderUid string, order models.Order) {
	cache.GetCache()[orderUid] = order
}

func GetFromCache(orderUid string) (models.Order, bool) {
	val, exists := cache.GetCache()[orderUid]
	return val, exists
}

func RestoreCache() {
	var orders []models.Order

	result := database.GetDB().Preload("Items").Preload("Payment").Preload("Delivery").Find(&orders)

	if result.Error != nil {
		log.Fatal("Can not load cache from db")
	}

	for _, order := range orders {
		cache.GetCache()[order.OrderUid] = order
	}
}
