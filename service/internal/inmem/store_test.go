package inmem_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

func TestAddAndGet(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	store := catcher.NewStoreService(inmem.NewStore())

	key := "TestAddAndGet"
	address := fmt.Sprintf("%s@%s", key, catcher.DefaultDomain)

	err := store.Add(key, catcher.Email{To: address})
	is.NoErr(err)

	emails := store.Get(key)
	is.Equal(1, emails.Len())

	email, ok := emails.At(0)
	is.True(ok)
	is.Equal(address, email.To)
}

func TestAddAndGetConcurrent(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	store := catcher.NewStoreService(inmem.NewStore())

	key := "TestAddAndGetConcurrent"
	address := fmt.Sprintf("%s@%s", key, catcher.DefaultDomain)
	initial := "initial@test.mx.ax"

	// add an initial email
	err := store.Add(key, catcher.Email{
		To:   address,
		From: initial,
	})
	is.NoErr(err)

	var total int32 = 1000
	var count int32

	// read constantly
	go func() {
		for atomic.LoadInt32(&count) < total {
			emails := store.Get(key)
			is.True(emails.Len() > 0)
			email, ok := emails.At(0)
			is.True(ok)
			is.Equal(initial, email.From)
		}
	}()

	var i int32
	for i = 0; i < total; i++ {
		err := store.Add(key, catcher.Email{
			To: address,
		})
		is.NoErr(err)
		atomic.AddInt32(&count, 1)
	}

	emails := store.Get(key)
	is.Equal(int(total)+1, emails.Len())
}

func TestSubscribeUnsubscribe(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	store := catcher.NewStoreService(inmem.NewStore())

	key := "TestSubscribeUnsubscribe"
	address := fmt.Sprintf("%s@%s", key, catcher.DefaultDomain)

	sub := store.Subscribe(key)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		email := <-sub.C
		is.Equal(address, email.To)
	}()

	time.Sleep(time.Second)

	err := store.Add(key, catcher.Email{
		To: address,
	})
	is.NoErr(err)

	wg.Wait()

	store.Unsubscribe(key, sub)
}
