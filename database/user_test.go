package database

import (
	"fmt"
	"testing"

	"github.com/ChrisMcKenzie/dropship/model"
)

func TestUserCreate(t *testing.T) {
	user := NewUser(
		"github.com",
		"chrismckenzie",
		"Chris McKenzie",
	)

	CreateUser(user)

	if user.Id == 0 {
		t.Fail()
	}

	db.Where(&model.User{Login: "chrismckenzie"}).First(&user)

	fmt.Println(user)

	if user.Id == 0 {
		t.Fail()
	}
}

func TestUserFindOrCreate(t *testing.T) {
	user := NewUser(
		"github.com",
		"chrismckenzie1",
		"Chris McKenzie",
	)

	FindOrCreateUser(user)

	if user.Id == 0 {
		t.Fail()
	}

	db.Where(&model.User{Login: "chrismckenzie"}).First(&user)

	fmt.Println(user)

	if user.Id == 0 {
		t.Fail()
	}
}
