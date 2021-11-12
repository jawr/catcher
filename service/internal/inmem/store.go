package inmem

import (
	"sync"

	"github.com/jawr/catcher/service/internal/catcher"
)

type store struct {
	storage    map[string]catcher.Emails
	storageMtx sync.RWMutex

	subscribers    map[string][]*catcher.Subscription
	subscribersMtx sync.RWMutex
}

// NewStore initialises a new inmemory store
func NewStore() *store {
	s := store{
		storage:     make(map[string]catcher.Emails),
		subscribers: make(map[string][]*catcher.Subscription),
	}
	return &s
}

// Add attempts to add an email to the correct address in the store
func (s *store) Add(key string, email catcher.Email) error {
	s.storageMtx.RLock()
	emails, ok := s.storage[key]
	if !ok {
		emails = catcher.NewEmails()
	}
	s.storageMtx.RUnlock()

	emails = emails.AddEmail(email)

	s.storageMtx.Lock()
	s.storage[key] = emails
	s.storageMtx.Unlock()

	// notify subscribers
	s.subscribersMtx.RLock()
	if subscribers, ok := s.subscribers[key]; ok {
		for _, sub := range subscribers {
			select {
			case sub.C <- emails:
			default:
			}
		}
	}
	s.subscribersMtx.RUnlock()

	return nil
}

// Has returns whether or not emails exist for a key
func (s *store) Has(key string) bool {
	s.storageMtx.RLock()
	defer s.storageMtx.RUnlock()

	_, ok := s.storage[key]

	return ok
}

// Get attempts to return all emails for a given address, it takes an immutable copy
func (s *store) Get(key string) catcher.Emails {
	s.storageMtx.RLock()
	defer s.storageMtx.RUnlock()

	emails, ok := s.storage[key]
	if !ok {
		return catcher.NewEmails()
	}

	return emails
}

// Subscribe lets a caller be notified whenever a new email for provided key
// is received
func (s *store) Subscribe(key string) *catcher.Subscription {
	subscription := &catcher.Subscription{
		C: make(chan catcher.Emails),
	}

	s.subscribersMtx.Lock()
	defer s.subscribersMtx.Unlock()

	subscribers, ok := s.subscribers[key]
	if !ok {
		subscribers = make([]*catcher.Subscription, 0)
	}

	subscribers = append(subscribers, subscription)

	s.subscribers[key] = subscribers

	return subscription
}

func (s *store) Unsubscribe(key string, a *catcher.Subscription) {
	s.subscribersMtx.Lock()
	defer s.subscribersMtx.Unlock()

	subscribers, ok := s.subscribers[key]
	if !ok {
		return
	}

	var i int
	for _, b := range subscribers {
		if a == b {
			continue
		}
		subscribers[i] = b
		i++
	}
	s.subscribers[key] = subscribers[:i]
	a = nil
}
