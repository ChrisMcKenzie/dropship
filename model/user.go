package model

import "time"

type User struct {
	Id          int64      `json:"id" gorm:"primary_key"`
	DeletedAt   *time.Time `json:"-"`
	CreatedAt   time.Time  `json:"created"`
	UpdatedAt   time.Time  `json:"updated"`
	Courier     string     `json:"courier"`
	Login       string     `json:"login"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Token       string     `json:"-"`
	TokenExpiry time.Time  `json:"-"`
}
