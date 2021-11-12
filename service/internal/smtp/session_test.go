package smtp

import (
	"bytes"
	"compress/gzip"
	"errors"
	"testing"

	"github.com/emersion/go-smtp"
	"github.com/google/gofuzz"
	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/matryer/is"
)

func TestSessionAuthPlain(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	f := fuzz.New()

	var username, password string

	f.Fuzz(&username)
	f.Fuzz(&password)

	var s session

	err := s.AuthPlain(username, password)
	is.NoErr(err)
}

func TestSessionMail(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	f := fuzz.New()

	var from string

	f.Fuzz(&from)

	var s session

	err := s.Mail(from, smtp.MailOptions{})
	is.NoErr(err)

	s.current.From = from
}


func TestSessionRcpt(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	cases := []struct{
		name string
		to string
		domain string
		err error
	}{
			{
				name: "invalid to",
				to: "to",
				err: catcher.ErrInvalid,
			},
			{
				name: "invalid domain",
				to: "to@example.com",
				domain: catcher.DefaultDomain,
				err: catcher.ErrInvalid,
			},
			{
				name: "valid domain",
				to: "to@example.com",
				domain: "example.com",
			},
		}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			s := session{
				domain: tc.domain,
			}
			err := s.Rcpt(tc.to)
			is.True(errors.Is(err, tc.err))
		})
	}
}

func TestSessionData(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	f := fuzz.New().NilChance(0)

	var data []byte

	f.Fuzz(&data)

	var s session

	buffer := bytes.NewBuffer(data)

	err := s.Data(buffer)
	is.NoErr(err)

	// decompress and validate
	decompressor, err := gzip.NewReader(bytes.NewBuffer(s.current.Data))
	is.NoErr(err)

	var readBuffer bytes.Buffer
	_, err = readBuffer.ReadFrom(decompressor)
	is.NoErr(err)

	is.Equal(data, readBuffer.Bytes())
}

func TestSessionReset(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	f := fuzz.New().NilChance(0)

	var email catcher.Email

	f.Fuzz(&email)

	s := session{current: email}

	s.Reset()

	is.Equal(0, len(s.current.To))
	is.Equal(0, len(s.current.From))
	is.Equal(0, len(s.current.Data))
}
