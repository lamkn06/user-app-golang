package response

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message"`
}
