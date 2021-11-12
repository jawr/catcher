package smtp

import (
	"fmt"
	"strings"
	"sync"

	"github.com/emersion/go-smtp"
	"github.com/jawr/catcher/service/internal/catcher"
)

// backend implements the smtp Backend interface
type backend struct {
	pool *sync.Pool
	handler catcher.EmailHandlerFn
}

// newBackend sets up and returns a new backend, sets up the session to push itself
// back into the pool after use
func newBackend(domain string, handler catcher.EmailHandlerFn) (*backend, error) {

	if len(domain) == 0 {
		return nil, fmt.Errorf("%w: no domain provided", catcher.ErrInvalid)
	}

	pool := &sync.Pool{}

	pool.New = func() interface{}{
		s := &session{
			emails: catcher.NewEmails(),
			domain: strings.ToLower(domain),
		}

		// session.onLogout is run on every logout/connection close
		// and is used to clean up the session object for reuse before putting
		// it back in to the pool
		s.onLogout = func() {
			for i := 0; i < s.emails.Len(); i++ {
				if email, ok := s.emails.At(i); ok {
					handler(email)
				}
			}

			// clear current
			s.Reset()

			// clear emails
			s.emails = catcher.NewEmails()

			// push back in to the queue for reuse
			pool.Put(s)
		}

		return s
	}

	b := backend{
		pool: pool,
	}

	return &b, nil
}

// Login returns smtp.ErrAuthUnsupported as we do not support it
func (b *backend) Login(_ *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, smtp.ErrAuthUnsupported
}

// AnonymousLogin returns nil and a new session as its allowed
func (b *backend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	s := b.pool.Get().(*session)
	return s, nil
}
