package catcher_test

import (
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/matryer/is"
)

func TestEmailValid(t *testing.T) {
	is := is.New(t)
	email := catcher.Email{
		From: "from",
		To: "to",
		Data: []byte("data"),
	}
	is.True(email.Valid())
}

func TestEmails(t *testing.T) {
	is := is.New(t)

	emails := catcher.NewEmails()
	is.Equal(0, emails.Len())

	expected := catcher.Email{
		From: "TestEmailsAddEmail",
	}

	next := emails.AddEmail(expected)
	is.Equal(0, emails.Len())
	is.Equal(1, next.Len())

	got, ok := next.At(0)
	is.True(ok)
	is.Equal(expected, got)

	got, ok = emails.At(0)
	is.True(!ok)
	is.Equal(catcher.Email{}, got)
}
