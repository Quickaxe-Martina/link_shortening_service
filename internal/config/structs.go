package config

// Config variables
type Config struct {
	RunAddr        string `env:"SERVER_ADDRESS"`
	ServerAddr     string `env:"BASE_URL"`
	DataFilePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn    string `env:"DATABASE_DSN"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
	DevMode        bool   `env:"DEV_MODE"`
	SecretKey      string `env:"SECRET_KEY"`
	TokenExp       int    `env:"TOKEN_EXP"`
}

// NewConfig create Config
func NewConfig() *Config {
	var cfg = Config{
		RunAddr:        "",
		ServerAddr:     "",
		DataFilePath:   "",
		DatabaseDsn:    "",
		MigrationsPath: "./migrations",
		DevMode:        false,
		SecretKey:      "my_secret_key",
		TokenExp:       3,
	}
	LoadEnv(&cfg)
	ParseFlags(&cfg, true)
	return &cfg
}
