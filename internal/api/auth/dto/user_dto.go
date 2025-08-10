package dto

type (
	RegisterRequest struct {
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	RegisterResponse struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
	}

	UserResponse struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
	}
)
