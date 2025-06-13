package dto

import (
	"github.com/Rafaelse6/mottus-ops-desafio/internal/entity"
)

type MotoDTO struct {
	Year  int    `json:"year"`
	Model string `json:"model"`
	Plate string `json:"plate"`
}

func NewMotoDTO(year int, model string, plate string) *MotoDTO {
	return &MotoDTO{
		Year:  year,
		Model: model,
		Plate: plate,
	}
}

func (d *MotoDTO) ToEntity() (*entity.Moto, error) {
	return entity.NewMoto(d.Year, d.Model, d.Plate)
}
