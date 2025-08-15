package config

// Config variables
type Config struct {
	RunAddr    string
	ServerAddr string
	URLData    map[string]string
}

// NewConfig create Config
func NewConfig(runAddr, serverAddr string) *Config {
	return &Config{
		RunAddr:    runAddr,
		ServerAddr: serverAddr,
		URLData:    make(map[string]string),
	}
}
