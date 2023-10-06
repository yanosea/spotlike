/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package config

type Config struct {
	cache *Cache
	cred  *Credential
}

func New() (*Config, error) {
	cache, err := LoadCache()
	if err != nil {
		return nil, err
	}

	cred, err := LoadCred()
	if err != nil {
		return nil, err
	}

	return &Config{
		cache: cache,
		cred:  cred,
	}, nil
}
