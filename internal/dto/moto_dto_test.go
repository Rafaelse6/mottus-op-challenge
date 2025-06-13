package dto

func NewMotoDTO(year int, model string, plate string) *MotoDTO {
	return &MotoDTO{
		Year:  year,
		Model: model,
		Plate: plate,
	}
}
func (m *MotoDTO) ToEntity() (*Moto, error) {
	moto, err := NewMoto(m.Year, m.Model, m.Plate)
	if err != nil {
		return nil, err
	}
	return moto, nil
}
