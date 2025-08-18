package dto_response

type (
	UserResponse struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
	}
)
