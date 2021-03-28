package logger

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	log := NewSTDLogger()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(log.CreateMiddleware())

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "?qwe=asd")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
