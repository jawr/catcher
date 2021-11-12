package inmem

import (
	"sync"
	"sync/atomic"

	"github.com/jawr/catcher/service/internal/catcher"
)

// producer is a very naive inmemory implementation that uses channels
type producer struct {
	queue chan <- catcher.Email 

	stopMtx sync.RWMutex
	stopped int32
}

// NewProducer accepts a channel and configures a new inmemory producer
func NewProducer(queue chan <- catcher.Email) *producer {
	return &producer{
		queue: queue,
	}
}

// Push blocks until the email can be pushed on to the queue
func (p *producer) Push(email catcher.Email) error {
	p.stopMtx.RLock()
	defer p.stopMtx.RUnlock()

	if atomic.LoadInt32(&p.stopped) > 0{
		return catcher.ErrProducerStopped
	}

	p.queue <- email

	return nil
}

// Stop prevents a producer from pushing any more messages. it locks stopMtx to ensure that when 
// it returns that the underlying queue can be closed
func (p *producer) Stop() error {
	p.stopMtx.Lock()
	defer p.stopMtx.Unlock()

	atomic.StoreInt32(&p.stopped, 1)

	return nil
}
