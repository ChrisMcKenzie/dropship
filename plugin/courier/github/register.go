package github

import (
	"github.com/ChrisMcKenzie/dropship/plugin/courier"
	"github.com/spf13/viper"
)

var (
	GithubClientId = viper.GetString("github.clientid")
	GithubSecret   = viper.GetString("github.secret")
)

func init() {
	courier.Register(
		New(GithubClientId, GithubSecret),
	)
}
