package dto

type CreateMotoRequest struct {
	Year  int    `json:"year" binding:"required"`
	Model string `json:"model" binding:"required"`
	Plate string `json:"plate" binding:"required"`
}

type MotoResponse struct {
	ID    string `json:"id"`
	Year  int    `json:"year"`
	Model string `json:"model"`
	Plate string `json:"plate"`
}
