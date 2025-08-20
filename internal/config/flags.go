package config

import (
	"flag"
	"strings"
)

// ParseFlags parses command-line flags to configure the server's runtime parameters.
func ParseFlags() *Config {
	var runAddr string
	var serverAddr = "http://localhost:8080/"

	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.Func("b", "server address before short URL", func(s string) error {
		if len(s) == 0 {
			return nil
		}
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		serverAddr = s
		return nil
	})
	flag.Parse()
	cfg := NewConfig(runAddr, serverAddr)
	return cfg
}
