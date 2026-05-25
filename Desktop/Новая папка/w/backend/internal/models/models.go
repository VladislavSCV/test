package models

import "gorm.io/gorm"
const StatusNew = "Новый"

type User struct {
	gorm.Model
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Login        string `json:"login" gorm:"uniqueIndex"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
}

type FoodOrder struct {
	gorm.Model
	UserID uint `json:"user_id"`
	User   User `json:"user,omitempty"`
	Address string `json:"address"`
	Dish string `json:"dish"`
	DeliveryTime string `json:"delivery_time"`
	Total int `json:"total"`
	Status string `json:"status"`
	Review *Review `json:"review,omitempty" gorm:"foreignKey:RecordID;references:ID"`
}
type Review struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	RecordID uint   `json:"record_id" gorm:"uniqueIndex"`
	Text     string `json:"text"`
	Rating   int    `json:"rating"`
}
