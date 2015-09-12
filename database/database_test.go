package database

import "testing"

func TestStoreTokenFor(t *testing.T) {
	err := StoreTokenFor("ChrisMcKenzie/dropship", "<token>")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestGetTokenFor(t *testing.T) {
	token, err := GetTokenFor("ChrisMcKenzie/dropship")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if token != "<token>" {
		t.Fail()
	}
}
