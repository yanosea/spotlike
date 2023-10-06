/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package config

import (
	"path/filepath"

	// https://github.com/spf13/viper
	"github.com/spf13/viper"
)

type Cache struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	cache Cache
)

func LoadCache() (*Cache, error) {
	viper.AddConfigPath("$XDG_CACHE_HOME/spotlike/")
	viper.SetConfigType("json")
	viper.SetConfigName("cache")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.WriteConfigAs(filepath.Join("$XDG_CACHE_HOME/spotlike/", "cache.json")); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func SaveCache(cache *Cache) error {
	var err error
	if err = viper.WriteConfigAs(filepath.Join("$XDG_CACHE_HOME/spotlike/", "cache.json")); err != nil {
		return err
	}
	return nil
}
