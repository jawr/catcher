package http

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	randomKeyLength = 10
	randomAttempts  = 10
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randomKey generates a random string of given length, it currently does not use numbers
// as locla parts of emails can not begin with them
func randomKey(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomEmailKeyResponse wraps a generated random key
type RandomEmailKeyResponse struct {
	Key string `json:"key"`
}

// handleRandomEmail attempts to return a key that currently has no emails attached
func (s *Server) handleRandomEmail(response http.ResponseWriter, request *http.Request) {
	for i := 0; i < randomAttempts; i++ {
		key := randomKey(randomKeyLength)

		if s.store.Has(key) {
			continue
		}

		r := RandomEmailKeyResponse{Key: key}

		if err := writeJSON(response, &r); err != nil {
			log.Printf("unable to encode key response for %q: %s", key, err)
			http.Error(response, "unable to encode key", http.StatusInternalServerError)
			return
		}

		return
	}

	// nothing found
	http.Error(response, "unable to generate random key", http.StatusInternalServerError)
}

// handleSubscribe accepts a websocket and attempts to send an initial set of emails before
// listening and sending more data. it also maintains the socket using a ping/pong
func (s *Server) handleSubscribe() http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"binary"},
		CheckOrigin: func(request *http.Request) bool {
			host := strings.ToLower(strings.Split(request.URL.Host, ":")[0])
			expected := strings.ToLower(strings.Split(s.httpd.Addr, ":")[0])
			return host == expected
		},
	}

	const (
		websocketPongWait   = 60 * time.Second
		websocketPingPeriod = (websocketPongWait * 9) / 10
		websocketWriteWait  = 100 * time.Second
	)

	return func(response http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(response, request, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Printf("unable to upgrade websocket: %s", err)
			}
			return
		}
		defer ws.Close()

		var subscriptionReq struct {
			Key string `json:"key"`
		}

		if err := ws.ReadJSON(&subscriptionReq); err != nil {
			log.Printf("unable to read subscription request: %s", err)
			return
		}

		subscription := s.store.Subscribe(subscriptionReq.Key)
		defer s.store.Unsubscribe(subscriptionReq.Key, subscription)

		// get initial state
		emails := s.store.Get(subscriptionReq.Key)
		if err := ws.WriteJSON(&emails); err != nil {
			log.Printf("unable to write initial emails data for %q: %s", subscriptionReq.Key, err)
			return
		}

		// writer
		go func() {
			ticker := time.NewTicker(websocketPingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-subscription.C:
					emails := s.store.Get(subscriptionReq.Key)
					if err := ws.WriteJSON(&emails); err != nil {
						log.Printf("unable to write emails to websocket for %q: %s", subscriptionReq.Key, err)
						return
					}

				case <-ticker.C:
					ws.SetWriteDeadline(time.Now().Add(websocketWriteWait))
					if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
						return
					}
				}
			}
		}()

		// reader
		ws.SetReadLimit(512)
		ws.SetReadDeadline(time.Now().Add(websocketPongWait))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(websocketPongWait))
			return nil
		})

		// noop any message sent down the line, breaking if we get any issues
		for {
			if err := request.Context().Err(); err != nil {
				break
			}
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}
