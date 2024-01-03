package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/entity"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServerRouter(t *testing.T) {
	st := repository.NewServerMemStorage()
	got := NewServerRouter(st)
	require.NotNil(t, got)
}

func TestMainPage(t *testing.T) {
	t.Run("Test root page", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestUpdateMetric(t *testing.T) {
	t.Run("Positive test", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/update/counter/someMetric/527", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Negative test with another method", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/update/counter/someMetric/527", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})
	t.Run("Negative with invalid data", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/update/counter/invalidTest/test", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetMetric(t *testing.T) {
	t.Run("Test get gauge metric", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/value/gauge/someMetric", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("Test get counter metric", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		request, err := http.NewRequest(http.MethodGet, ts.URL+"/value/counter/someMetric", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestUpdateMetricJSON(t *testing.T) {
	t.Run("Test update counter metric", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		for i := 0; i < 5; i++ {
			temp := int64(i)
			m := entity.Metrics{
				MType: "counter",
				ID:    "test",
				Delta: &temp,
			}
			body, err := json.Marshal(m)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/update/", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			client := &http.Client{}
			resp, err := client.Do(request)
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			fmt.Println(string(b))
			defer resp.Body.Close()
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
	t.Run("Test update gauge metric", func(t *testing.T) {
		controller := repository.NewServerMemStorage()
		r := NewServerRouter(controller)
		ts := httptest.NewServer(r)
		defer ts.Close()
		for i := 0; i < 5; i++ {
			temp := float64(i)
			m := entity.Metrics{
				MType: "gauge",
				ID:    "test",
				Value: &temp,
			}
			body, err := json.Marshal(m)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/update/", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			client := &http.Client{}
			resp, err := client.Do(request)
			require.NoError(t, err)
			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			fmt.Println(string(b))
			defer resp.Body.Close()
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
}
