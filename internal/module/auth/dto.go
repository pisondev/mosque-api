package auth

type LoginGoogleRequest struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	Role        string `json:"role"`
}

type AccountProfileResponse struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	AvatarURL       string `json:"avatar_url,omitempty"`
	Role            string `json:"role"`
	GoogleConnected bool   `json:"google_connected"`
}

type UpdateAccountProfileRequest struct {
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}
