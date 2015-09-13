package database

import (
	"fmt"
	"testing"
)

func TestUserCreate(t *testing.T) {
	user := CreateUser(
		"github.com",
		"chrismckenzie",
		"Chris McKenzie",
		"chris@chrismckenzie.io",
		"bd2d40d2d399280d23bcecbccb74c117",
	)

	if user.Id == 0 {
		t.Fail()
	}

	db.Where(&User{Login: "chrismckenzie"}).First(&user)

	fmt.Println(user)

	if user.Id == 0 {
		t.Fail()
	}
}
