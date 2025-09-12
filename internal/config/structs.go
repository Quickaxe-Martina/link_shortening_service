package config

// Config variables
type Config struct {
	RunAddr      string `env:"SERVER_ADDRESS"`
	ServerAddr   string `env:"BASE_URL"`
	DataFilePath string `env:"FILE_STORAGE_PATH"`
	URLData      map[string]string
}

// NewConfig create Config
func NewConfig() *Config {
	var cfg = Config{
		RunAddr:    "",
		ServerAddr: "",
		URLData:    make(map[string]string),
	}
	LoadEnv(&cfg)
	ParseFlags(&cfg, true)
	return &cfg
}
