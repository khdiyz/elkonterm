package models

type LoginRequest struct {
	Email    string `json:"email" binding:"required" default:"admin@mail.ru"`
	Password string `json:"password" binding:"required" default:"admin"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
