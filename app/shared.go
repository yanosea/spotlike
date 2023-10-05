/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package app

type Shared struct {
	cache *Cache
	cred  *Credential
}

func New() *Shared {
	return &Shared{
		cache: nil,
		cred:  nil,
	}
}
