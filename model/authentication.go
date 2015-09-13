package model

import "time"

type Authentication struct {
	Login  string
	Email  string
	Token  string
	Expiry time.Time
}
