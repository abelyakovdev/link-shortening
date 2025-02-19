package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestURL     string
		expectedStatus int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "valid key",
			requestURL:     "",
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "<a href=\"http://example.com\">Temporary Redirect</a>.\n\n",
			expectedHeader: "http://example.com",
		},
		{
			name:           "invalid key",
			requestURL:     "/invalidkey",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not found\n",
			expectedHeader: "",
		},
		{
			name:           "incorrect parameters",
			requestURL:     "/incorrect/incorrect",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Incorrect parameters\n",
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedHeader != "" {
				storage = make(map[string]string)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.expectedHeader))
				w := httptest.NewRecorder()

				PostHandler(w, req)

				resp := w.Result()
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)

				req = httptest.NewRequest(http.MethodGet, string(body), nil)
				w = httptest.NewRecorder()

				GetHandler(w, req)

				resp = w.Result()
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
				assert.Equal(t, tt.expectedBody, string(body))
				assert.Equal(t, tt.expectedHeader, resp.Header.Get("Location"))
			} else {
				req := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)
				w := httptest.NewRecorder()

				GetHandler(w, req)

				resp := w.Result()
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)

				require.NoError(t, err)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedPrefix string
		expectedURL    string
	}{
		{
			name:           "valid URL",
			requestBody:    "http://example.com",
			expectedStatus: http.StatusCreated,
			expectedPrefix: "http://localhost:8080/",
			expectedURL:    "http://example.com",
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedPrefix: "",
			expectedURL:    "Incorrect parameters\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage = make(map[string]string)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			PostHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			if tt.expectedPrefix != "" {
				responseString := string(body)
				assert.True(t, strings.HasPrefix(responseString, tt.expectedPrefix))

				key := strings.TrimPrefix(responseString, tt.expectedPrefix)
				storedURL, exists := storage[key]
				require.True(t, exists)
				assert.Equal(t, tt.expectedURL, storedURL)
			} else {
				assert.Equal(t, tt.expectedURL, string(body))
			}
		})
	}
}
