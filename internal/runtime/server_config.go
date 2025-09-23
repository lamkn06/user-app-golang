package runtime

type ServerConfig struct {
	Port        string `env:"PORT" envDefault:"8081"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	APIVersion  string `env:"API_VERSION" envDefault:"v1"`
}
