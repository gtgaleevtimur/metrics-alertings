package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAgentMemStorage(t *testing.T) {
	ms := NewAgentMemStorage()
	require.NotNil(t, ms)
}

func TestUpdateMemStorage(t *testing.T) {
	ms := NewAgentMemStorage()
	require.NotNil(t, ms)
	ms.UpdateMemStorage()
}

func TestSendMetrics(t *testing.T) {
	type request struct {
		datatype string
		key      any
		value    any
	}

	// http server response body
	response := "response"

	tests := []struct {
		name     string
		req      request
		wanterr  error
		wantcode int
	}{
		{
			name: "Test Valid Post request gauge metric",
			req: request{
				datatype: "gauge",
				key:      "m01",
				value:    1.34,
			},
			wanterr:  nil,
			wantcode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				http.Error(rw, response, tc.wantcode)
			}))
			defer server.Close()

			req := server.URL + "/update/" + tc.req.datatype + fmt.Sprintf("/%v/%v", tc.req.key, tc.req.value)
			ms := NewAgentMemStorage()
			err := ms.SendMetrics(req)
			assert.Equal(t, tc.wanterr, err)
		})
	}
}
