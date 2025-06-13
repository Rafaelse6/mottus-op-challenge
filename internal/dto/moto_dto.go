package dto

type MotoDTO struct {
	Year  int    `json:"year"`
	Model string `json:"model"`
	Plate string `json:"plate"`
}

func ToEntity(motoDTO *MotoDTO) (*Moto, error) {
	moto, err := NewMoto(motoDTO.Year, motoDTO.Model, motoDTO.Plate)
	if err != nil {
		return nil, err
	}
	return moto, nil
}

func NewMoto(i int, s1, s2 string) (any, any) {
	panic("unimplemented")
}
