package response

type SignUpResponse struct {
	User NewUserResponse `json:"user"`
}

type SignInResponse struct {
	Token        string          `json:"token"`
	RefreshToken string          `json:"refresh_token"`
	User         NewUserResponse `json:"user"`
}

type SignOutResponse struct {
	Message string `json:"message"`
}
