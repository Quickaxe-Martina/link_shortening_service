package config

import (
	"flag"
	"strings"
)

// ParseFlags parses command-line flags to configure the server's runtime parameters.
func ParseFlags(cfg *Config) {
	var serverAddr = "http://localhost:8080/"

	flag.StringVar(&cfg.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.DataFilePath, "f", "", "saved urls path")
	flag.StringVar(&cfg.DatabaseDsn, "d", "", "database")
	useTLS := flag.Bool("s", false, "")
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
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&cfg.AuditURL, "audit-url", "", "remote audit receiver URL")
	flag.Parse()
	cfg.ServerAddr = serverAddr
	cfg.UseTLS = *useTLS
}
