package config

// apiConfig stores the configuration for the API
// DefaultLimitPagination is the default limit for pagination
// MaxLimitPagination is the maximum limit for pagination
type apiConfig struct {
	DefaultLimitPagination int
	MaxLimitPagination     int
}

var apiConfigInstance apiConfig

// GetApiConfig returns the apiConfig instance
func GetApiConfig() apiConfig {
	if apiConfigInstance.DefaultLimitPagination == 0 {
		apiConfigInstance = apiConfig{
			DefaultLimitPagination: 20,
			MaxLimitPagination:     200,
		}
	}

	return apiConfigInstance
}
