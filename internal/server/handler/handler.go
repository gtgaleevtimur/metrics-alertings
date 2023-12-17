package handler

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository"
	"net/http"
	"strconv"
)

func NewServerRouter(repository repository.ServerStorager) *chi.Mux {
	controller := newServerHandler(repository)
	router := chi.NewRouter()
	router.Get("/", controller.MainPage)
	router.Post("/update/{type}/{metric}/{value}", controller.UpdateMetric)
	router.Get("/value/gauge/{metric}", controller.GetMetric)
	router.Get("/value/counter/{metric}", controller.GetMetric)
	router.MethodNotAllowedHandler()
	router.NotFoundHandler()
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	return router
}

type ServerHandler struct {
	Repository repository.ServerStorager
}

func newServerHandler(repository repository.ServerStorager) *ServerHandler {
	return &ServerHandler{
		Repository: repository,
	}
}

func (h *ServerHandler) MainPage(res http.ResponseWriter, req *http.Request) {
	body := `
<!DOCTYPE html>
<html>
    <head>
        <title>All tuples</title>
    </head>
    <body>
      <table>
          <tr>
            <td>Metric</td>
            <td>Value</td>
          </tr>
    `
	list := h.Repository.List()
	for k, v := range list {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}
	body = body + " </table>\n </body>\n</html>"

	res.Write([]byte(body))
}

func (h *ServerHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	memType := chi.URLParam(req, "type")
	metric := chi.URLParam(req, "metric")
	value := chi.URLParam(req, "value")

	switch memType {
	case "counter":
		v, err := strconv.Atoi(value)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		err = h.Repository.Update(memType, metric, int64(v))
		if err != nil {
			http.Error(res, "http.StatusBadRequest", http.StatusBadRequest)
		}
		res.WriteHeader(http.StatusOK)
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		err = h.Repository.Update(memType, metric, v)
		if err != nil {
			http.Error(res, "http.StatusBadRequest", http.StatusBadRequest)
		}
		res.WriteHeader(http.StatusOK)
	default:
		http.Error(res, "Metric type is not valid", http.StatusBadRequest)
	}
}

func (h *ServerHandler) GetMetric(res http.ResponseWriter, req *http.Request) {
	metric := chi.URLParam(req, "metric")
	v, err := h.Repository.Get(metric)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("%v", v)))
}
