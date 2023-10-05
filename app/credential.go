/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Credential struct {
	Client Client `toml:"client"`
}

type Client struct {
	ClientId     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

func LoadConfig() (*Credential, error) {
	viper.AddConfigPath("$XDG_CONFIG_HOME/spotlike/")
	viper.SetConfigName("credential")
	viper.SetConfigType("toml")

	viper.SetDefault("client_id", "")
	viper.SetDefault("client_secret", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.WriteConfigAs(filepath.Join("$XDG_CONFIG_HOME/spotlike/", "credential.toml")); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var cred Credential
	if err := viper.Unmarshal(&cred); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &cred, nil
}
