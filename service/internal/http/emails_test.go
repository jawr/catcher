package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

var testConfig = Config{
	Address: "0.0.0.0:0",
}

func readExampleEmail(is *is.I) []byte {
	is.Helper()
	data, err := ioutil.ReadFile("testdata/example.email.gz")
	is.NoErr(err)
	return data
}

func TestHandleEmailsNoEmail(t *testing.T) {
	is := is.New(t)

	request := httptest.NewRequest(http.MethodGet, "/emails/", nil)
	recorder := httptest.NewRecorder()

	store := catcher.NewStoreService(inmem.NewStore())

	server := NewServer(testConfig, store)

	server.handleEmails(recorder, request)

	is.Equal(http.StatusBadRequest, recorder.Code)
}

func TestHandleEmailsNotFound(t *testing.T) {
	is := is.New(t)

	request := httptest.NewRequest(http.MethodGet, "/empty", nil)

	request = mux.SetURLVars(request, map[string]string{
		"key": "foobar@bar.com",
	})
	recorder := httptest.NewRecorder()

	store := catcher.NewStoreService(inmem.NewStore())

	server := NewServer(testConfig, store)

	server.handleEmails(recorder, request)

	is.Equal(http.StatusOK, recorder.Code)

	response := recorder.Result()
	defer response.Body.Close()

	var emails []catcher.Email
	err := json.NewDecoder(response.Body).Decode(&emails)
	is.NoErr(err)

	is.Equal(0, len(emails))
}

func TestHandleEmails(t *testing.T) {
	is := is.New(t)

	request := httptest.NewRequest(http.MethodGet, "/empty", nil)
	request = mux.SetURLVars(request, map[string]string{
		"key": "foobar",
	})
	recorder := httptest.NewRecorder()

	store := catcher.NewStoreService(inmem.NewStore())

	expected := 10

	for i := 0; i < expected; i++ {
		store.Add("foobar", catcher.Email{
			To:   "foobar@bar.com",
			Data: readExampleEmail(is),
		})
	}

	server := NewServer(testConfig, store)

	server.handleEmails(recorder, request)

	is.Equal(http.StatusOK, recorder.Code)

	response := recorder.Result()
	defer response.Body.Close()

	var emails []catcher.Email
	err := json.NewDecoder(response.Body).Decode(&emails)
	is.NoErr(err)

	is.Equal(expected, len(emails))
}

func TestHandleRandomEmail(t *testing.T) {
	is := is.New(t)

	request := httptest.NewRequest(http.MethodGet, "/random", nil)
	recorder := httptest.NewRecorder()

	store := catcher.NewStoreService(inmem.NewStore())

	server := NewServer(testConfig, store)

	server.handleRandomEmail(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	var key RandomEmailKeyResponse
	err := json.NewDecoder(response.Body).Decode(&key)
	is.NoErr(err)

	is.Equal(randomKeyLength, len(key.Key))
}
