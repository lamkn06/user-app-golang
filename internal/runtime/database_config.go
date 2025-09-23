package runtime

import "fmt"

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"local"`
	Password string `env:"DB_PASSWORD" envDefault:"local"`
	DBName   string `env:"DB_NAME" envDefault:"db_name"`
}

func (c *DatabaseConfig) PrimaryConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DBName)
}
