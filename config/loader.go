package config

import (
	"github.com/argcv/webeh/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"github.com/pkg/errors"
	"fmt"
)

var (
	kProjectName = "manul"
)

// this function will search and load configurations
func LoadConfig(path string) (err error) {
	viper.SetConfigName(kProjectName)
	viper.SetEnvPrefix(kProjectName)

	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", kProjectName))
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", kProjectName))
		if conf := os.Getenv(fmt.Sprintf("%s_CFG", strings.ToUpper(kProjectName))); conf != "" {
			viper.SetConfigFile(conf)
		}
	}

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
		log.Errorf("Load configure failed: %s", err.Error())
		return err
	}
	if conf := viper.ConfigFileUsed(); conf != "" {
		wd, _ := os.Getwd()
		if rel, _ := filepath.Rel(wd, conf); rel != "" && strings.Count(rel, "..") < 3 {
			conf = rel
		}
		log.Infof("Using config file: %s", conf)
		return nil
	} else {
		msg := "No configure file"
		log.Warnf(msg)
		return errors.New(msg)
	}
}