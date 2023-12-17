package repository

type AgentStorager interface {
	UpdateMemStorage()
	SendMetrics(addr string) error
}
