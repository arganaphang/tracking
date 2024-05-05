package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Data    struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type ActiveUserResponse struct {
	Message string `json:"message"`
}
