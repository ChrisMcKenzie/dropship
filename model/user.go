package model

import "time"

type User struct {
	Id          int64     `gorm:"primary_key" json:"-"`
	Courier     string    `json:"json"`
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	AccessToken string    `json:"-"`
	Gravatar    string    `json:"gravatar"`
	CreatedAt   time.Time `json:"created"`
	UpdatedAt   time.Time `json:"updated"`
}
