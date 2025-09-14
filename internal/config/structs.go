package config

// Config variables
type Config struct {
	RunAddr      string `env:"SERVER_ADDRESS"`
	ServerAddr   string `env:"BASE_URL"`
	DataFilePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn  string `env:"DATABASE_DSN"`
}

// NewConfig create Config
func NewConfig() *Config {
	var cfg = Config{
		RunAddr:    "",
		ServerAddr: "",
	}
	LoadEnv(&cfg)
	ParseFlags(&cfg, true)
	return &cfg
}
