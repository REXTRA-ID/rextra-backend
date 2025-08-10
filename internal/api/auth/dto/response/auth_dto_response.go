package authDto_response

type (
	RegisterResponse struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
	}

	LoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}

	LoginWithGoogleResponse struct {
		NeedRegistration bool   `json:"need_registration"`
		Token            string `json:"token"`
		RegisterToken    string `json:"register_token"`
		Role             string `json:"role"`
	}
)
