package repository

type ServerStorager interface {
	Update(memType, metric string, value interface{}) error
	Get(metric string) (interface{}, error)
	List() map[string]interface{}
}
