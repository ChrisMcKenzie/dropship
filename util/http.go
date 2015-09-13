package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

func GenerateRandom() string {
	size := 32 // change the length of the generated random string here

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	return rs
}

func CreateCookie(w http.ResponseWriter, name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}

func GetCookieValue(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	return c.Value, err
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "nil",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
