package repository

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
)

var contentType = url.Values{"Content-type": {"text/plain"}}

type AgentMemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewAgentMemStorage() AgentStorager {
	return &AgentMemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
	}
}

func (m *AgentMemStorage) UpdateMemStorage() {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	m.Gauge["Alloc"] = float64(stat.Alloc)
	m.Gauge["BuckHashSys"] = float64(stat.BuckHashSys)
	m.Gauge["Frees"] = float64(stat.Frees)
	m.Gauge["GCCPUFraction"] = float64(stat.GCCPUFraction)
	m.Gauge["GCSys"] = float64(stat.GCSys)
	m.Gauge["HeapAlloc"] = float64(stat.HeapAlloc)
	m.Gauge["HeapIdle"] = float64(stat.HeapIdle)
	m.Gauge["HeapInuse"] = float64(stat.HeapInuse)
	m.Gauge["HeapObjects"] = float64(stat.HeapObjects)
	m.Gauge["HeapReleased"] = float64(stat.HeapReleased)
	m.Gauge["HeapSys"] = float64(stat.HeapSys)
	m.Gauge["LastGC"] = float64(stat.LastGC)
	m.Gauge["Lookups"] = float64(stat.Lookups)
	m.Gauge["MCacheInuse"] = float64(stat.MCacheInuse)
	m.Gauge["MCacheSys"] = float64(stat.MCacheSys)
	m.Gauge["MSpanInuse"] = float64(stat.MSpanInuse)
	m.Gauge["MSpanSys"] = float64(stat.MSpanSys)
	m.Gauge["Mallocs"] = float64(stat.Mallocs)
	m.Gauge["NextGC"] = float64(stat.NextGC)
	m.Gauge["NumForcedGC"] = float64(stat.NumForcedGC)
	m.Gauge["NumGC"] = float64(stat.NumGC)
	m.Gauge["OtherSys"] = float64(stat.OtherSys)
	m.Gauge["PauseTotalNs"] = float64(stat.PauseTotalNs)
	m.Gauge["StackInuse"] = float64(stat.StackInuse)
	m.Gauge["StackSys"] = float64(stat.StackSys)
	m.Gauge["Sys"] = float64(stat.Sys)
	m.Gauge["TotalAlloc"] = float64(stat.TotalAlloc)
	m.Gauge["RandomValue"] = rand.Float64()
	m.Counter["PollCount"] += 1
}

func (m *AgentMemStorage) SendMetrics(addr string) error {
	for k, v := range m.Gauge {
		req := fmt.Sprintf("%sgauge/%v/%v", addr, k, v)
		res, err := http.PostForm(req, contentType)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			line, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("%s: %s; %s",
				"Can't send report to the server",
				res.Status,
				line)
		}
		res.Body.Close()
	}
	for k, v := range m.Counter {
		req := fmt.Sprintf("%scounter/%v/%v", addr, k, v)
		res, err := http.PostForm(req, contentType)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			line, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("%s: %s; %s",
				"Can't send report to the server",
				res.Status,
				line)
		}
		res.Body.Close()
	}
	return nil
}
