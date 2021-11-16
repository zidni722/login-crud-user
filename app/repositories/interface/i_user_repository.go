package _interface

import (
	"github.com/jinzhu/gorm"
	"github.com/zidni722/login-crud-user/app/repositories"
)

type IUserRepository interface {
	repositories.BaseRepository
	CreateUser(db *gorm.DB, entity interface{}) error
	FindByUsername(db *gorm.DB, entity interface{}, username string) error
}