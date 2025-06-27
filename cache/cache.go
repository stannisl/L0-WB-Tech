package cache

import (
	"L0-wb/models"
)

var cacheMap = make(map[string]models.Order)

func GetCache() map[string]models.Order {
	return cacheMap
}
