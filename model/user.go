package model

import "time"

type User struct {
	Id          int64     `gorm:"primary_key" json:"-"`
	Courier     string    `json:"json"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Token       string    `json:"-"`
	TokenExpiry time.Time `json:"-"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
}
