package repository

import (
	"fmt"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/entity"
)

type ServerMemStorage struct {
	Data map[string]interface{}
}

func NewServerMemStorage() ServerStorager {
	return &ServerMemStorage{
		Data: make(map[string]interface{}),
	}
}

func (m *ServerMemStorage) List() map[string]interface{} {
	return m.Data
}

func (m *ServerMemStorage) Get(metric string) (interface{}, error) {
	if v, ok := m.Data[metric]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *ServerMemStorage) Update(memType, metric string, value interface{}) error {
	switch memType {
	case "counter":
		if m.Data[metric] == nil {
			m.Data[metric] = value
			return nil
		}
		m.Data[metric] = m.Data[metric].(int64) + value.(int64)
		return nil
	case "gauge":
		m.Data[metric] = value
		return nil
	default:
		return fmt.Errorf("wrong type of metric")
	}
}

func (m *ServerMemStorage) UpdateJSON(metric *entity.Metrics) (*entity.Metrics, error) {
	result := &entity.Metrics{
		MType: metric.MType,
		ID:    metric.ID,
	}
	switch metric.MType {
	case "counter":
		if m.Data[metric.ID] == nil {
			m.Data[metric.ID] = *metric.Delta
			val := m.Data[metric.ID].(int64)
			result.Delta = &val
			return result, nil
		}
		m.Data[metric.ID] = m.Data[metric.ID].(int64) + *metric.Delta
		val := m.Data[metric.ID].(int64)
		result.Delta = &val
		return result, nil
	case "gauge":
		m.Data[metric.ID] = *metric.Value
		val := m.Data[metric.ID].(float64)
		result.Value = &val
		return result, nil
	default:
		return nil, fmt.Errorf("wrong type of metric")
	}
}

func (m *ServerMemStorage) GetJSON(metric *entity.Metrics) (*entity.Metrics, error) {
	result := &entity.Metrics{
		MType: metric.MType,
		ID:    metric.ID,
	}
	switch metric.MType {
	case "counter":
		if val, ok := m.Data[metric.ID]; ok {
			tempVal := val.(int64)
			result.Delta = &tempVal
			return result, nil
		}
		return nil, entity.ErrNoFound
	case "gauge":
		if val, ok := m.Data[metric.ID]; ok {
			tempVal := val.(float64)
			result.Value = &tempVal
			return result, nil
		}
		return nil, entity.ErrNoFound
	default:
		return nil, fmt.Errorf("wrong type of metric")
	}
}
