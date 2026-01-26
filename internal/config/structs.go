/*
Package config for config
*/
package config

// Config variables
// generate:reset
type Config struct {
	RunAddr            string   `env:"SERVER_ADDRESS" json:"server_address"`
	ServerAddr         string   `env:"BASE_URL" json:"base_url"`
	DataFilePath       string   `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDsn        string   `env:"DATABASE_DSN" json:"database_dsn"`
	MigrationsPath     string   `env:"MIGRATIONS_PATH" json:"migrations_path"`
	DevMode            bool     `env:"DEV_MODE" json:"dev_mode"`
	SecretKey          string   `env:"SECRET_KEY" json:"secret_key"`
	TokenExp           int      `env:"TOKEN_EXP" json:"token_exp"`
	DeleteBatchSize    int      `env:"DELETE_BATCH_SIZE" json:"delete_batch_size"`
	DeleteTimeDuration int      `env:"DELETE_TIME_DURATION" json:"delete_time_duration"`
	AuditFile          string   `env:"AUDIT_FILE" json:"audit_file"`
	AuditURL           string   `env:"AUDIT_URL" json:"audit_url"`
	UseTLS             bool     `env:"ENABLE_HTTPS" json:"enable_https"`
	ShutdownTimeout    int      `env:"SHUTDOWN_TIMEOUT" json:"shutdown_timeout"`
	HostWhitelist      []string `env:"HOST_WHITELIST" envSeparator:"," json:"host_whitelist"`
	TrustedSubnet      string   `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCAddr           string   `env:"GRPC_ADDR" json:"grpc_addr"`
	ConfigPath         string
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
		DeleteBatchSize:    50,
		AuditFile:          "./audit_data.json",
		AuditURL:           "",
		UseTLS:             false,
		ShutdownTimeout:    10,
		HostWhitelist:      nil,
		TrustedSubnet:      "",
		GRPCAddr:           "",
	}
	ParseFlags(&cfg)
	LoadConfigFile(&cfg)
	LoadEnv(&cfg)
	return &cfg
}
