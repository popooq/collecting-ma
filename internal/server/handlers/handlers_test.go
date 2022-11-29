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

/** func TestMetricStorageServeHTTP(t *testing.T) {

	tests := []struct {
		name string
		url  string
		code int
	}{
		{
			name: "Positive test: Gauge",
			url:  "/update/gauge/Alloc/12",
			code: 200,
		},
		{
			name: "Positive test: Counter",
			url:  "/update/counter/PollCount/111",
			code: 200,
		},
		{
			name: "Negative test: Unkonwn metric",
			url:  "/update/unknown/poop/111",
			code: 501,
		},
		{
			name: "Negative test: Empty value",
			url:  "/update/counter/Name1",
			code: 404,
		},
		{
			name: "Negative test: Wrong counter value",
			url:  "/update/counter/name1/sdf",
			code: 400,
		},
		{
			name: "Negative test: Wrong gauge value",
			url:  "/update/gauge/name3/sdf",
			code: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			memStorage := storage.NewMemStorage()
			handler := NewmetricStorage(memStorage)

			h := http.Handler(handler)
			h.ServeHTTP(w, r)
			result := w.Result()
			if result.StatusCode != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, result.StatusCode)
			}
			defer result.Body.Close()
		})
	}
} **/

func NewRouter() chi.Router {
	memS := storage.NewMemStorage()
	handler := NewmetricStorage(memS)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/value/{mType}/{mName}", func(w http.ResponseWriter, r *http.Request) {
			handler.MetricValue(w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.AllMetrics(w, r)
		})
		r.Post("/update/{mType}/{mName}/{mValue}", func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		})
	})
	return r
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Client()

	statusCode, body := testRequest(t, ts, "POST", "/update/gauge/name/12")
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

	statusCode, _ = testRequest(t, ts, "GET", "/value/gauge/name")
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
