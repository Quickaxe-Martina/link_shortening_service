/*
Package config for config
*/
package config

// Config variables
// generate:reset
type Config struct {
	RunAddr            string `env:"SERVER_ADDRESS"`
	ServerAddr         string `env:"BASE_URL"`
	DataFilePath       string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn        string `env:"DATABASE_DSN"`
	MigrationsPath     string `env:"MIGRATIONS_PATH"`
	DevMode            bool   `env:"DEV_MODE"`
	SecretKey          string `env:"SECRET_KEY"`
	TokenExp           int    `env:"TOKEN_EXP"`
	DeleteBachSize     int    `env:"DELETE_BACH_SIZE"`
	DeleteTimeDuration int    `env:"DELETE_TIME_DURATION"`
	AuditFile          string `env:"AUDIT_FILE"`
	AuditURL           string `env:"AUDIT_URL"`
}

// NewConfig create Config
func NewConfig() *Config {
	var cfg = Config{
		RunAddr:            "",
		ServerAddr:         "",
		DataFilePath:       "",
		DatabaseDsn:        "",
		MigrationsPath:     "./migrations",
		DevMode:            false,
		SecretKey:          "my_secret_key",
		TokenExp:           3,
		DeleteTimeDuration: 5,
		DeleteBachSize:     50,
		AuditFile:          "./audit_data.json",
		AuditURL:           "",
	}
	LoadEnv(&cfg)
	ParseFlags(&cfg, true)
	return &cfg
}
