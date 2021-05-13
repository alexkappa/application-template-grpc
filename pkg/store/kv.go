package store

import (
	"context"
	"errors"
)

var (
	ErrKeyNotFount = errors.New("key not found")
)

// KVStore is a simple key-value data store abstraction.
type KVStore interface {
	Get(ctx context.Context, key string) (value interface{}, err error)
	Set(ctx context.Context, key string, value interface{}) (err error)
}

// NewInMemoryKVStore is an implementation of KVStore that doesn't persist
// any data on disk. It's created for illustration purposes
func NewInMemoryKVStore() KVStore {
	return &inMemoryKVStore{make(map[string]interface{})}
}

type inMemoryKVStore struct {
	m map[string]interface{}
}

func (s inMemoryKVStore) Get(ctx context.Context, key string) (interface{}, error) {
	if v, ok := s.m[key]; ok {
		return v, nil
	}
	return nil, ErrKeyNotFount
}

func (s inMemoryKVStore) Set(ctx context.Context, key string, value interface{}) error {
	s.m[key] = value
	return nil
}
