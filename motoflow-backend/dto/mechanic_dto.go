package dto

type CreateMechanicRequest struct {
	Code           string `json:"code" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Phone          string `json:"phone"`
	Specialization string `json:"specialization"`
	PitNumber      int    `json:"pit_number"`
	IsAvailable    bool   `json:"is_available"`
}

type UpdateMechanicRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Specialization string `json:"specialization"`
	PitNumber      int    `json:"pit_number"`
	IsAvailable    *bool  `json:"is_available"`
}

