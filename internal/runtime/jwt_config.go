package runtime

import "time"

type JWTConfig struct {
	SecretKey     string        `env:"JWT_SECRET_KEY" envDefault:"your-secret-key"`
	Expiration    time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRATION" envDefault:"168h"` // 7 days
}
