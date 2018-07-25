package config

import (
	"github.com/spf13/viper"
)

func checkKeyIsExists(key string) bool {
	return viper.Get(key) != nil
}

func getStringOrDefault(key, or string) string {
	if checkKeyIsExists(key) {
		return viper.GetString(key)
	} else {
		return or
	}
}

func getStringSliceOrDefault(key string, or []string) []string {
	if checkKeyIsExists(key) {
		return viper.GetStringSlice(key)
	} else {
		return or
	}
}

func getIntOrDefault(key string, or int) int {
	if checkKeyIsExists(key) {
		return viper.GetInt(key)
	} else {
		return or
	}
}

func getInt64OrDefault(key string, or int64) int64 {
	if checkKeyIsExists(key) {
		return viper.GetInt64(key)
	} else {
		return or
	}
}

func getBoolOrDefault(key string, or bool) bool {
	if checkKeyIsExists(key) {
		return viper.GetBool(key)
	} else {
		return or
	}
}

func setConfig(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}
