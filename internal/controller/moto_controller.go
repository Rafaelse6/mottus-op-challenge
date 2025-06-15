package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/dto"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MotoController struct {
	service *service.MotoService
}

func NewMotoController(service *service.MotoService) *MotoController {
	return &MotoController{service: service}
}

func (c *MotoController) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.MotoDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	moto, err := c.service.CreateMoto(input.Year, input.Model, input.Plate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(moto)
}

func (c *MotoController) List(w http.ResponseWriter, r *http.Request) {
	plate := r.URL.Query().Get("plate")
	motos, err := c.service.ListMotos(plate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(motos)
}

func (c *MotoController) UpdatePlate(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body struct {
		Plate string `json:"plate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdatePlate(id, body.Plate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *MotoController) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteMoto(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
