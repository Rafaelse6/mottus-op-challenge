package entity

import (
	"errors"
	"strings"

	"github.com/Rafaelse6/mottus-ops-desafio/pkg/entity"
)

type Moto struct {
	ID    entity.ID `json:"id"`
	Year  int       `json:"year"`
	Model string    `json:"model"`
	Plate string    `json:"plate"`
}

func NewMoto(year int, model string, plate string) (*Moto, error) {

	if year <= 0 {
		return nil, errors.New("year must be positive")
	}

	if strings.TrimSpace(plate) == "" {
		return nil, errors.New("plate cannot be empty")
	}

	return &Moto{
		ID:    entity.NewID(),
		Year:  year,
		Model: model,
		Plate: plate,
	}, nil
}
