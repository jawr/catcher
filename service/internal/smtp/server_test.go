package smtp_test

import (
	"errors"
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/smtp"
	"github.com/matryer/is"
)

func noopEmailHandlerFn(_ catcher.Email) error { return nil }

func TestServerInvalidConfig(t *testing.T) {
	is := is.New(t)

	cases := []struct {
		name   string
		domain string
		config smtp.Config
	}{
		{
			name: "empty domain",
			config: smtp.Config{
				Addr: "foo",
			},
		},
		{
			name:   "empty addr",
			domain: catcher.DefaultDomain,
			config: smtp.Config{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			_, err := smtp.NewServer(tc.domain, tc.config, noopEmailHandlerFn)
			is.True(errors.Is(err, catcher.ErrInvalid))
		})
	}
}
