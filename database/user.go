package database

import "github.com/ChrisMcKenzie/dropship/model"

func CreateUser(courier, login, name, email, gravatar string) *model.User {
	user := new(model.User)
	user.Courier = courier
	user.Login = login
	user.Name = name
	user.Email = email
	user.Gravatar = gravatar

	db.Create(user)
	return user
}
