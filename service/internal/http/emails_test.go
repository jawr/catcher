package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/matryer/is"
)

var testConfig = Config{
	Address:  "0.0.0.0:0",
	SiteRoot: "testdata/index.html",
}

func readExampleEmail(is *is.I) []byte {
	is.Helper()
	data, err := ioutil.ReadFile("testdata/example.email.gz")
	is.NoErr(err)
	return data
}

func TestHandleRandomEmail(t *testing.T) {
	is := is.New(t)

	request := httptest.NewRequest(http.MethodGet, "/random", nil)
	recorder := httptest.NewRecorder()

	store := catcher.NewStoreService(inmem.NewStore())

	server, err := NewServer(testConfig, nil, store)
	is.NoErr(err)

	server.handleRandomEmail(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	var key RandomEmailKeyResponse
	err = json.NewDecoder(response.Body).Decode(&key)
	is.NoErr(err)

	is.Equal(randomKeyLength, len(key.Key))
}

// TODO: websocket handler tests
