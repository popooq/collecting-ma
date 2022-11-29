package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/popooq/collectimg-ma/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewRouter() chi.Router {
	memS := storage.NewMemStorage()
	handler := NewmetricStorage(memS)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/update/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		})
		r.Get("/value/{mType}/{mName}", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricValue(w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
	})
	return r
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Client()

	statusCode, body := testRequest(t, ts, "POST", "/update/gauge/name/123")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "", body)

	statusCode, body = testRequest(t, ts, "POST", "/update/counter/name/123")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "", body)

	statusCode, _ = testRequest(t, ts, "POST", "/update/gauge/name/asdf")
	assert.Equal(t, http.StatusBadRequest, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/name/dasd")
	assert.Equal(t, http.StatusBadRequest, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/sdfa/name/12")
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = testRequest(t, ts, "GET", "/")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "GET", "/value/guage/name")
	assert.Equal(t, http.StatusNotImplemented, statusCode)

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()
	return resp.StatusCode, string(respBody)
}
