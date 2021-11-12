package catcher

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/textproto"
	"time"

	"github.com/jhillyerd/enmime"
)

type emailJSON struct {
	From       string               `json:"from"`
	To         string               `json:"to"`
	Subject    string               `json:"subject"`
	Headers    textproto.MIMEHeader `json:"headers"`
	HTML       []byte               `json:"html"`
	Text       []byte               `json:"text"`
	ReceivedAt time.Time            `json:"received_at"`
}

func (e Emails) MarshalJSON() ([]byte, error) {
	emails := make([]emailJSON, 0)

	for _, e := range e.emails {
		reader, err := gzip.NewReader(bytes.NewBuffer(e.Data))
		if err != nil {
			return nil, fmt.Errorf("unable to decompress email data: %w", err)
		}

		envelope, err := enmime.ReadEnvelope(reader)
		if err != nil {
			return nil, fmt.Errorf("unable to read envelope: %w", err)
		}

		emails = append(emails, emailJSON{
			From:       e.From,
			To:         e.To,
			Subject:    envelope.GetHeader("subject"),
			Headers:    envelope.Root.Header,
			HTML:       []byte(envelope.HTML),
			Text:       []byte(envelope.Text),
			ReceivedAt: e.ReceivedAt,
		})

	}

	return json.Marshal(&emails)
}
