/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package app

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Cache struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	cache Cache
	err   error = nil
)

func LoadCache() (*Cache, error) {
	viper.AddConfigPath("$XDG_CACHE_HOME/spotlike/")
	viper.SetConfigType("json")
	viper.SetConfigName("cache")

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.WriteConfigAs(filepath.Join("$XDG_CACHE_HOME/spotlike/", "cache.json")); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	if err != nil {
		if err = viper.Unmarshal(&cache); err != nil {
			fmt.Println(err)
		}
	}

	return &cache, err
}

func SaveCache(cache *Cache) error {
	if err = viper.WriteConfigAs(filepath.Join("$XDG_CACHE_HOME/spotlike/", "cache.json")); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
