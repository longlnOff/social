package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/longlnOff/social/internal/auth"
	"github.com/longlnOff/social/internal/store"
	"github.com/longlnOff/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	testAuth := &auth.TestAuthenticator{}

	logger := zap.NewNop()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStore:    mockCacheStore,
		authenticator: testAuth,
	}
}


func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("WANT %d BUT GOT %d", expected, actual)
	}
}
