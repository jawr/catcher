package catcher

import (
	"strings"
)

// StoreService wraps a store and provides some generic transformations, namely normalising the key
type StoreService struct {
	store Store
}

// NewStoreService wraps an existing store and provides some generic transformations
func NewStoreService(store Store) *StoreService {
	return &StoreService{
		store: store,
	}
}

// Add wraps the underlying store feeding in a transformed key
func (s *StoreService) Add(key string, email Email) error {
	// normalise email key
	parts := strings.Split(key, "@")
	key = strings.ToLower(parts[0])
	return s.store.Add(key, email)
}

// Get wraps the underlying store feeding in a transformed key
func (s *StoreService) Get(key string) Emails {
	return s.store.Get(strings.ToLower(key))
}

// Has wraps the underlying store feeding in a transformed key
func (s *StoreService) Has(key string) bool {
	return s.store.Has(strings.ToLower(key))
}

// Subscribe wraps the underlying store feeding in a transformed key
func (s *StoreService) Subscribe(key string) *Subscription {
	return s.store.Subscribe(strings.ToLower(key))
}

// Unsubscribe wraps the underlying store feeding in a transformed key
func (s *StoreService) Unsubscribe(key string, sub *Subscription) {
	s.store.Unsubscribe(strings.ToLower(key), sub)
}
