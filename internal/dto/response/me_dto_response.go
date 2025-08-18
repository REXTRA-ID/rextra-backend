package dto_response

type (
	GetMe struct {
		PersonalInfo PersonalInfo `json:"personal_info"`
	}

	PersonalInfo struct {
		ID          string `json:"id"`
		Fullname    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
	}
)
