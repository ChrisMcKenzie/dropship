package model

import "time"

type Repo struct {
	Id        int64     `gorm:"primary_key" json:"-"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Host      string    `json:"host"`
	Courier   string    `json:"courier"`
	URL       string    `json:"url"`
	CloneURL  string    `json:"clone_url"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}
