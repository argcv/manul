package config

import (
	"fmt"
	"strings"
)

type SMTPAuth struct {
	Username string
	Password string
}

func GetSMTPAuth() *SMTPAuth {
	if getBoolOrDefault(KeyMailSMTPPerformAuth, false) {
		// with auth
		user := getStringOrDefault(KeyMailSMTPUserName, "")
		pass := getStringOrDefault(KeyMailSMTPPassword, "")

		return &SMTPAuth{
			Username: user,
			Password: pass,
		}
	}
	return nil
}

type SMTPConfig struct {
	Host               string
	Port               int
	Sender             string
	DefaultFrom        string
	InsecureSkipVerify bool
	Auth               *SMTPAuth
}

func (c *SMTPConfig) GetUsername() string {
	if c.Auth == nil {
		return ""
	} else {
		return c.Auth.Username
	}
}

func (c *SMTPConfig) GetPassword() string {
	if c.Auth == nil {
		return ""
	} else {
		return c.Auth.Password
	}
}

func GetSMTPConfig() (cfg *SMTPConfig) {
	cfg = &SMTPConfig{
		Host:               getStringOrDefault(KeyMailSMTPHost, ""),
		Port:               getIntOrDefault(KeyMailSMTPPort, 0),
		Sender:             getStringOrDefault(KeyMailSMTPUserSender, ""),
		InsecureSkipVerify: getBoolOrDefault(KeyMailSMTPInsecureSkipVerify, false),
		Auth:               GetSMTPAuth(),
	}
	if cfg.Auth != nil && cfg.DefaultFrom == "" {
		if strings.Contains(cfg.Auth.Username, "@") {
			cfg.DefaultFrom = cfg.Auth.Username
		} else {
			cfg.DefaultFrom = fmt.Sprintf("%s@%s", cfg.Auth.Username, cfg.Host)
		}
	}
	return
}
