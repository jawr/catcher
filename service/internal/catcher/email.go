package catcher

import (
	"time"
)

// Email encapsulates all required data for an email
type Email struct {
	From string
	// To is the test.mx.ax email address
	To string
	// Data encapsulates the entire email message
	Data []byte

	ReceivedAt time.Time
}

// Valid does a rudimentary check to see if the email is valid
func (e *Email) Valid() bool {
	return len(e.From) > 0 && len(e.To) > 0 && len(e.Data) > 0
}

// Emails is a wrapper around a slice of emails that allows us to do less copies
type Emails struct {
	emails []Email
}

func NewEmails() Emails {
	return Emails{
		emails: make([]Email, 0),
	}
}

// AddEmail allows us to safely add an email to our slice
func (e Emails) AddEmail(email Email) Emails {
	if e.emails == nil {
		e = NewEmails()
	}
	e.emails = append(e.emails, email)
	return e
}

// Len returns the number of elements in Email
func (e Emails) Len() int {
	return len(e.emails)
}

// At returns the email at index
func (e Emails) At(idx int) (Email, bool) {
	if len(e.emails) < idx+1 {
		return Email{}, false
	}
	return e.emails[idx], true
}
