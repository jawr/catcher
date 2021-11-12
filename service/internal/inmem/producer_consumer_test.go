package inmem_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

func TestProducerConsumer(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	queue := make(chan catcher.Email)
	producer := inmem.NewProducer(queue)
	consumer := inmem.NewConsumer(queue)

	expected := 10
	got := 0

	var wg sync.WaitGroup
	wg.Add(expected)

	go consumer.Handler(func (_ catcher.Email) error {
		defer wg.Done()
		got++
		return nil
	})

	for i := 0; i < expected; i++ {
		err := producer.Push(catcher.Email{})
		is.NoErr(err)
	}
	err := producer.Stop()
	is.NoErr(err)

	wg.Wait()

	is.Equal(expected, got)

	err = producer.Push(catcher.Email{})
	is.True(errors.Is(err, catcher.ErrProducerStopped))
}
