package memory

import "fmt"

type ServerMemStorage struct {
	Data map[string]interface{}
}

func NewServerMemStorage() *ServerMemStorage {
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
	case "gauge":
		m.Data[metric] = value
	default:
		return fmt.Errorf("wrong type of metric")
	}
	return nil
}
