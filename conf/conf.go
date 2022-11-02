package conf

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%v", err)
	}
}

func Conf() *viper.Viper {
	return viper.GetViper()
}

func CaptchaCachePrefix() string {
	prefix := Conf().GetString("CAPTCHA_CACHE_PREFIX")
	if prefix == "" {
		return "captcha:"
	}
	return prefix
}

func TokenCachePrefix() string {
	cache_prefix := Conf().GetString("TOKEN_CACHE_PREFIX")
	if cache_prefix == "" {
		return "token:"
	}
	return cache_prefix
}
