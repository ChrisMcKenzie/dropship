package httputil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
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

// GetScheme is a helper function that evaluates the http.Request
// and returns the scheme, HTTP or HTTPS. It is able to detect,
// using the X-Forwarded-Proto, if the original request was HTTPS
// and routed through a reverse proxy with SSL termination.
func GetScheme(r *http.Request) string {
	switch {
	case r.URL.Scheme == "https":
		return "https"
	case r.TLS != nil:
		return "https"
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return "https"
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return "https"
	default:
		return "http"
	}
}

// GetHost is a helper function that evaluates the http.Request
// and returns the hostname. It is able to detect, using the
// X-Forarded-For header, the original hostname when routed
// through a reverse proxy.
func GetHost(r *http.Request) string {
	switch {
	case len(r.Host) != 0:
		return r.Host
	case len(r.URL.Host) != 0:
		return r.URL.Host
	case len(r.Header.Get("X-Forwarded-For")) != 0:
		return r.Header.Get("X-Forwarded-For")
	case len(r.Header.Get("X-Host")) != 0:
		return r.Header.Get("X-Host")
	case len(r.Header.Get("XFF")) != 0:
		return r.Header.Get("XFF")
	default:
		return "localhost:8080"
	}
}

func GetBaseURL(r *http.Request) string {
	return GetScheme(r) + "://" + GetHost(r)
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
