package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemiyKew/http-rest-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUserCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", nil)
	s := newServer(teststore.New())
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
