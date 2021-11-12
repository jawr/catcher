package smtp

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/jawr/catcher/service/internal/catcher"
)

// session represents a client connection to the server
// and implements github.com/emersion/go-smtp interface
type session struct {
	domain string

	// email state
	current catcher.Email
	emails  catcher.Emails

	// destructor
	onLogout func()
}

// AuthPlain handles any user authentication, which we do not need
func (s *session) AuthPlain(username, password string) error {
	return nil
}

// Mail is a handler for the MAIL FROM smtp command
func (s *session) Mail(from string, opts smtp.MailOptions) error {
	s.current.From = from
	return nil
}

// Rcpt is a handler for the RCPT smtp command
func (s *session) Rcpt(to string) error {
	parts := strings.Split(to, "@")
	if len(parts) != 2 {
		return fmt.Errorf("%w: email does not contain a '@'", catcher.ErrInvalid)
	}

	if strings.ToLower(parts[1]) != s.domain {
		return fmt.Errorf("%w: invalid recipient", catcher.ErrInvalid)
	}

	s.current.To = to
	return nil
}

// Data is a handler for the DATA smtp command
func (s *session) Data(data io.Reader) error {
	reader, writer := io.Pipe()

	// compress the data and pipe it in to our buffer data -> compressor -> reader
	go func() {
		defer writer.Close()

		compressor := gzip.NewWriter(writer)
		defer compressor.Close()

		io.Copy(compressor, data)
	}()

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(reader)
	if err != nil {
		log.Printf("unable to read data from %s to %s: %s", s.current.From, s.current.To, err)
		return errors.New("unable to read data")
	}

	s.current.Data = buffer.Bytes()
	s.current.ReceivedAt = time.Now()

	// push current in to the emails list and then continue, its the clients
	// duty to call reset if they want to continue with more emails
	s.emails = s.emails.AddEmail(s.current)

	return nil
}

// Reset is a handler for the RESET smtp command
func (s *session) Reset() {
	s.current.Data = nil
	s.current.From = ""
	s.current.To = ""
	s.current.ReceivedAt = time.Time{}
}

// Logout is a handler for the LOGOUT smtp command and is called when closing the client
// connection
func (s *session) Logout() error {
	s.onLogout()
	return nil
}
