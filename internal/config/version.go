package config

const (
	// APIVersion represents the current API version
	APIVersion = "v1"

	// APIPrefix is the base prefix for all API routes
	APIPrefix = "/api"
)

// GetVersionedAPIPath returns the versioned API path
func GetVersionedAPIPath(path string) string {
	return APIPrefix + "/" + APIVersion + path
}
