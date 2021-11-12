package inmem_test

import (
	"sync"
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

func TestConsumerHandler(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	queue := make(chan catcher.Email)
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
		queue <- catcher.Email{}
	}

	wg.Wait()

	close(queue)

	is.Equal(expected, got)
}
