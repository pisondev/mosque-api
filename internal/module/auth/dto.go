package auth

type LoginGoogleRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
