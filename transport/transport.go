package transport

import (
	"fmt"
	"net/http"

	"wb/domain"
	"wb/service"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type VacationResponse struct {
	VacationDays uint `json:"vacation_days"`
}

type SearchResponse struct {
	Employees []*domain.Employee `json:"employees"`
}

type API struct {
	service service.IService
	server  http.ServeMux
}

func New(service service.IService, port string) (*API, error) {
	api := &API{
		service: service,
	}

	r := mux.NewRouter()

	r.Use(Middleware())

	r.HandleFunc("/hire", api.hireEmployee).Methods(http.MethodPost)
	r.HandleFunc("/fire", api.fireEmployee).Methods(http.MethodDelete)
	r.HandleFunc("/vacation_days", api.getVacationDays).Methods(http.MethodGet)
	r.HandleFunc("/search", api.searchEmployee).Methods(http.MethodGet)

	portStr := fmt.Sprintf(":%s", port)

	log.Info().Msgf("Start HTTP server on %s port", port)

	if err := http.ListenAndServe(portStr, r); err != nil {
		return nil, err
	}

	return api, nil
}
