package models

import (
	"time"
)

type Order struct {
	Id          uint   `json:"-" gorm:"PrimaryKey"`
	OrderUid    string `json:"order_uid" gorm:"unique"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`

	DeliveryID uint     `json:"-"`
	Delivery   Delivery `json:"delivery" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	PaymentID uint    `json:"-"`
	Payment   Payment `json:"payment" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Items []Item `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
