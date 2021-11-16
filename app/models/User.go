package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username    string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Address string `json:"address" validate:"required"`
	Password  string `json:"password" validate:"required"`
}
