package github

import (
	"github.com/ChrisMcKenzie/dropship/plugin/courier"
	"github.com/spf13/viper"
)

func Register() {
	GithubClientId := viper.GetString("github.clientid")
	GithubSecret := viper.GetString("github.secret")
	courier.Register(
		New(GithubClientId, GithubSecret),
	)
}
