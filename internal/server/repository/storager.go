package repository

import "github.com/gtgaleevtimur/metrics-alertings/internal/server/entity"

type ServerStorager interface {
	Update(memType, metric string, value interface{}) error
	Get(metric string) (interface{}, error)
	List() map[string]interface{}
	UpdateJSON(metric *entity.Metrics) (*entity.Metrics, error)
	GetJSON(metric *entity.Metrics) (*entity.Metrics, error)
}
