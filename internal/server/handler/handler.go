package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/entity"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository"
	"net/http"
	"strconv"
)

const contentTypeJSON = "application/json"

func NewServerRouter(repository repository.ServerStorager) *chi.Mux {
	controller := newServerHandler(repository)
	router := chi.NewRouter()
	router.Use(LoggerHandler)
	router.Get("/", controller.MainPage)
	router.Post("/update/{type}/{metric}/{value}", controller.UpdateMetric)
	router.Get("/value/gauge/{metric}", controller.GetMetric)
	router.Get("/value/counter/{metric}", controller.GetMetric)

	router.Post("/update/", controller.UpdateMetricJSON)
	router.Post("/value/", controller.GetMetricJSON)

	router.MethodNotAllowedHandler()
	router.NotFoundHandler()
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

func (h *ServerHandler) UpdateMetricJSON(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != contentTypeJSON {
		http.Error(res, "wrong content type", http.StatusBadRequest)
	}
	var m entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	result, err := h.Repository.UpdateJSON(&m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", contentTypeJSON)
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (h *ServerHandler) GetMetricJSON(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != contentTypeJSON {
		http.Error(res, "wrong content type", http.StatusBadRequest)
	}
	var m entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	result, err := h.Repository.GetJSON(&m)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", contentTypeJSON)
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
