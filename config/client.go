package config

import (
	"github.com/argcv/go-argcvapis/app/manul/user"
)

func GetClientUserName() string {
	return getStringOrDefault(KeyClientUserName, "")
}

func GetClientUserSecret() string {
	return getStringOrDefault(KeyClientUserSecret, "")
}

func SetClientUserName(user string) error {
	return setConfig(KeyClientUserName, user)
}

func SetClientUserSecret(secret string) error {
	return setConfig(KeyClientUserSecret, secret)
}

func GetAuthInfo() (auth *user.AuthToken) {
	auth = &user.AuthToken{
		Name:   GetClientUserName(),
		Secret: GetClientUserSecret(),
	}
	return
}
