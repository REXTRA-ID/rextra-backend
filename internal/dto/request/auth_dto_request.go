package dto_request

type (
	RegisterRequest struct {
		Fullname    string `json:"fullname" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	LoginWithGoogleRequest struct {
		IdToken string `json:"id_token" binding:"required"`
	}

	LogoutRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	ChangePasswordRequest struct {
		NewPassword string `json:"new_password"`
	}
)
