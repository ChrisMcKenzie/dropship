package database

import "github.com/ChrisMcKenzie/dropship/model"

func NewUser(courier, login, name string) *model.User {
	user := new(model.User)
	user.Courier = courier
	user.Login = login
	user.Name = name
	return user
}

func CreateUser(u *model.User) {
	db.Create(u)
}

func GetUser(user *model.User) {
	db.Where(user).First(user)
}

func FindOrCreateUser(user *model.User) {
	db.Where(user).FirstOrCreate(user)
}

func UpdateUser(user *model.User) {
	db.Model(user).Update(user)
}
