package helpers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/stretchr/testify/assert"
)

func TestParseClientError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		action   string
		name     string
		err      error
		expected string
	}{
		"openapi": {
			action:   "create",
			name:     "radarr_tag",
			err:      &radarr.GenericOpenAPIError{},
			expected: "Unable to create radarr_tag, got error: \nDetails:\n",
		},
		"generic": {
			action:   "create",
			name:     "radarr_tag",
			err:      errors.New("other error"),
			expected: "Unable to create radarr_tag, got error: other error",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseClientError(test.action, test.name, test.err))
		})
	}
}

func TestParseNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind     string
		field    string
		search   string
		expected string
	}{
		"generic": {
			kind:     "radarr_tag",
			field:    "label",
			search:   "test",
			expected: "Unable to find radarr_tag, got error: data source not found: no radarr_tag with label 'test'",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseNotFoundError(test.kind, test.field, test.search))
		})
	}
}

func TestIsNotFound(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		status   int
		expected bool
	}{
		"not found": {
			status:   http.StatusNotFound,
			expected: true,
		},
		"unauthorized": {
			status:   http.StatusUnauthorized,
			expected: false,
		},
		"server error": {
			status:   http.StatusInternalServerError,
			expected: false,
		},
		"ok": {
			status:   http.StatusOK,
			expected: false,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte("{}"))
			}))
			defer server.Close()

			config := radarr.NewConfiguration()
			config.Servers = radarr.ServerConfigurations{{URL: server.URL}}
			_, _, err := radarr.NewAPIClient(config).TagAPI.GetTagById(context.Background(), 1).Execute()

			assert.Equal(t, test.expected, IsNotFound(err))
		})
	}
}

func TestIsNotFoundOtherErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		expected bool
	}{
		"nil": {
			err:      nil,
			expected: false,
		},
		"not an API error": {
			err:      errors.New("404 Not Found"),
			expected: false,
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, IsNotFound(test.err))
		})
	}
}
