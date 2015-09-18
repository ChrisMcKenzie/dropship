package model

import "time"

type Repo struct {
	Id        int64     `gorm:"primary_key" json:"-"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Courier   string    `json:"courier"`
	URL       string    `json:"url"`
	CloneURL  string    `json:"clone_url"`
	Active    bool      `json:"active"`
	Private   bool      `json:"-"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}
