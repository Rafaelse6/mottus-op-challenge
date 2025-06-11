package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafaelse6/mottus-ops-desafio/internal/dto"
	"github.com/Rafaelse6/mottus-ops-desafio/internal/service"
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
