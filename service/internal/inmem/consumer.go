package inmem

import (
	"log"

	"github.com/jawr/catcher/service/internal/catcher"
)

// consumer is a very
type consumer struct {
	queue <-chan catcher.Email
}

// NewConsumer creates a new inmemory consumer
func NewConsumer(queue <-chan catcher.Email) *consumer {
	return &consumer{
		queue: queue,
	}
}

// Handle accepts an email 
func (c *consumer) Handler(fn catcher.EmailHandlerFn) {
	for email := range c.queue {
		if err := fn(email); err != nil {
			log.Printf("unable to handle email: %s", err)
		}
	}
}
