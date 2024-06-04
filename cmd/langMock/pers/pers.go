package pers

import "nerd/yago1/cmd/langMock/store"

func Lookup(s store.Store, key string) ([]byte, error) {
	return s.Get(key)
}
