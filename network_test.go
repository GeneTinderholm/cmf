package cmf

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	t.Run("should use default options if not provided", func(t *testing.T) {
		wasCalled := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wasCalled = true
			assert.Equal(t, http.MethodGet, r.Method)
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Empty(t, body)
		}))
		defer ts.Close()

		_, _ = Fetch(ts.URL)

		assert.True(t, wasCalled)
	})

	t.Run("should use whatever method the user provides", func(t *testing.T) {
		wasCalled := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wasCalled = true
			assert.Equal(t, "SALAMI", r.Method)
		}))
		defer ts.Close()

		_, _ = Fetch(ts.URL, FetchOptions{Method: "SALAMI"})

		assert.True(t, wasCalled)
	})
	t.Run("should use whatever body the user provides", func(t *testing.T) {
		wasCalled := false
		expectedBody := "this is a body"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wasCalled = true
			bs, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, expectedBody, string(bs))
			_ = r.Body.Close()
		}))
		defer ts.Close()

		_, _ = Fetch(ts.URL, FetchOptions{Body: strings.NewReader(expectedBody)})

		assert.True(t, wasCalled)
	})
}

func TestFetchString(t *testing.T) {
	t.Run("should return the response as a string", func(t *testing.T) {
		expectedResponse := "this is a response"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(expectedResponse)); err != nil {
				t.Fail()
			}
		}))
		defer ts.Close()

		str, err := FetchString(ts.URL)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, str)
	})
}

func TestFetchJSON(t *testing.T) {
	type Thing struct {
		X int `json:"x"`
	}
	t.Run("should json decode the response", func(t *testing.T) {
		expectedResponse := `{"x": 1}`
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(expectedResponse)); err != nil {
				t.Fail()
			}
		}))
		defer ts.Close()

		thing, err := FetchJSON[Thing](ts.URL)
		assert.NoError(t, err)
		assert.Equal(t, Thing{X: 1}, thing)
	})
}
