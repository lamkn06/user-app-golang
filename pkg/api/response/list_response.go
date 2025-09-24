package response

// ListMeta represents metadata for list responses
type ListMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Meta    ListMeta `json:"meta"`
	Items   []T      `json:"items"`
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
}

// NewListResponse creates a new list response
func NewListResponse[T any](items []T, total int64, page, limit int) ListResponse[T] {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return ListResponse[T]{
		Meta: ListMeta{
			Total:       total,
			Page:        page,
			Limit:       limit,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
		Items:   items,
		Success: true,
	}
}
