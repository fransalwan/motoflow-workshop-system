package dto

type CreateCustomerRequest struct {
	Name         string `json:"name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	PlateNumber  string `json:"plate_number" binding:"required"`
	MotorModel   string `json:"motor_model" binding:"required"`
	Year         int    `json:"year"`
	FrameNumber  string `json:"frame_number"`
	EngineNumber string `json:"engine_number"`
	CurrentKM    int    `json:"current_km"`
}

type UpdateCustomerRequest struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	PlateNumber  string `json:"plate_number"`
	MotorModel   string `json:"motor_model"`
	Year         int    `json:"year"`
	FrameNumber  string `json:"frame_number"`
	EngineNumber string `json:"engine_number"`
	CurrentKM    int    `json:"current_km"`
}

