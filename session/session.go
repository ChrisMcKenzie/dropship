package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/ChrisMcKenzie/dropship/database"
	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
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

func CreateToken(audience string, user *model.User) (string, error) {
	secret := viper.GetString("secret")
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user_id"] = user.Id
	token.Claims["audience"] = audience
	token.Claims["expires"] = time.Now().UTC().Add(time.Hour * 72).Unix()
	return token.SignedString([]byte(secret))
}

func GetUserFromJWT(token string) (user *model.User) {
	secret := viper.GetString("secret")
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Panic(err)
		return nil
	}

	id, ok := t.Claims["user_id"].(float64)
	if !ok {
		return nil
	}

	user = &model.User{Id: int64(id)}
	database.GetUser(user)

	return user
}
