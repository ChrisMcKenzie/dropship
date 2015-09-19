package model

import "time"

type Model struct {
	Id        int        `json:"id" gorm:"primary_key"`
	DeletedAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
}
