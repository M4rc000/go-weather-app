package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `form:"Name" validate:"required"`
	Username string `gorm:"unique" form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required,min=6"`
}
