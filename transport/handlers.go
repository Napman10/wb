package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wb/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (a *API) hireEmployee(w http.ResponseWriter, r *http.Request) {
	var e domain.Employee

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	if _, err := a.service.HireEmployee(ctx, e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	log.Debug().Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id"))).Msg("success")
}

func (a *API) fireEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	ctx := r.Context()

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	if err = a.service.FireEmployee(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	log.Debug().Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id"))).Msg("success")
}

func (a *API) getVacationDays(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	ctx := r.Context()

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	vacationDays, err := a.service.GetVacationDays(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	resp := VacationResponse{VacationDays: vacationDays}

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	log.Debug().
		Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id"))).
		Str("response", string(respJson)).
		Msg("success")

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (a *API) searchEmployee(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	ctx := r.Context()

	employees, err := a.service.SearchEmployee(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	resp := SearchResponse{Employees: employees}

	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id")))

		return
	}

	log.Debug().
		Str("response", string(respJson)).
		Str("request_id", fmt.Sprintf("%v", ctx.Value("request-id"))).
		Msg("success")

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
