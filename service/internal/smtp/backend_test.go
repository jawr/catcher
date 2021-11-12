package smtp

import (
	"sync"
	"testing"

	"github.com/emersion/go-smtp"
	fuzz "github.com/google/gofuzz"
	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/matryer/is"
)

var testAllowedDomains = []string{"catcher.mx.ax"}

func noopEmailHandlerFn (_ catcher.Email) error {
	return nil
}

func TestBackendLogin(t *testing.T) {
	is := is.New(t)
	f := fuzz.New()
	var username, password string

	f.Fuzz(&username)
	f.Fuzz(&password)

	b, err := newBackend(catcher.DefaultDomain, noopEmailHandlerFn)
	is.NoErr(err)

	s, err := b.Login(nil, username, password)
	is.Equal(smtp.ErrAuthUnsupported, err)
	is.Equal(nil, s)
}

func TestBackendAnonymousLogin(t *testing.T) {
	is := is.New(t)

	b, err := newBackend(catcher.DefaultDomain, noopEmailHandlerFn)
	is.NoErr(err)

	_, err = b.AnonymousLogin(nil)
	is.NoErr(err)
}

func TestBackendEmailHandler(t *testing.T) {
	is := is.New(t)
	f := fuzz.New().Funcs(
		func(o *catcher.Email, c fuzz.Continue) {
			c.FuzzNoCustom(o)
			for o.From == "" { c.Fuzz(&o.From) }
			for o.To == "" { c.Fuzz(&o.To) }
			for len(o.Data) == 0 { c.Fuzz(&o.Data) }
		},
	)	

	var expected catcher.Email
	f.Fuzz(&expected)

	var wg sync.WaitGroup
	wg.Add(1)

	emailHandlerFn := func(email catcher.Email) error {
		defer wg.Done()
		is.Equal(expected, email)
		return nil
	}

	b, err := newBackend(catcher.DefaultDomain, emailHandlerFn)
	is.NoErr(err)

	s, err := b.AnonymousLogin(nil)	
	is.NoErr(err)

	s.(*session).emails = catcher.NewEmails().AddEmail(expected)

	err = s.Logout()
	is.NoErr(err)

	wg.Wait()
}
