package model

import "time"

type Authentication struct {
	Login  string
	Email  string
	Name   string
	Token  string
	Expiry time.Time
	Model
}
