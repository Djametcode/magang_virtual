package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
    Username string `json:"username" gorm:"unique;not null" validate:"required"`
    Email    string `json:"email" gorm:"unique;not null" validate:"required"`
    Password string `json:"password" gorm:"not null"`
    Photos   []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Photo struct {
    ID       uint      `gorm:"primaryKey;autoIncrement"`
    URL      string    `gorm:"not null" validate:"required"`
    UserID   uint      `gorm:"not null"`
    User     User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type LoginUser struct {
    gorm.Model
    Email    string `json:"email" gorm:"unique;not null" validate:"required"`
    Password string `json:"password" gorm:"not null"` 
}