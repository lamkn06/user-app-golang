package request

// ListRequest represents pagination parameters for list endpoints
type ListRequest struct {
	Page  int `json:"page" form:"page" query:"page" validate:"min=1"`
	Limit int `json:"limit" form:"limit" query:"limit" validate:"min=1,max=100"`
}

// NewListRequest creates a new ListRequest with default values
func NewListRequest() ListRequest {
	return ListRequest{
		Page:  1,
		Limit: 10,
	}
}

// GetPage returns the page number, defaulting to 1
func (r *ListRequest) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

// GetLimit returns the limit, defaulting to 10 and maxing at 100
func (r *ListRequest) GetLimit() int {
	if r.Limit <= 0 {
		return 10
	}
	if r.Limit > 100 {
		return 100
	}
	return r.Limit
}

// GetOffset calculates the offset for database queries
func (r *ListRequest) GetOffset() int {
	return (r.GetPage() - 1) * r.GetLimit()
}
