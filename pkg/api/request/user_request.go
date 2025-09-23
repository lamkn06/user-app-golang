package request

type NewUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
