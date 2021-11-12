package inmem_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	queue := make(chan catcher.Email)
	producer := inmem.NewProducer(queue)

	expected := 10
	got := 0

	go func() {
		for i := 0; i < expected; i++ {
			err := producer.Push(catcher.Email{})
			is.NoErr(err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(expected)

	go func() {
		for range queue {
			got++
			wg.Done()
		}
	}()

	wg.Wait()

	close(queue)

	is.Equal(expected, got)
}

func TestStopProducer(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	queue := make(chan catcher.Email)
	producer := inmem.NewProducer(queue)

	err := producer.Stop()
	is.NoErr(err)

	err = producer.Push(catcher.Email{})
	is.True(errors.Is(err, catcher.ErrProducerStopped))
}
