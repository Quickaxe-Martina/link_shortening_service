package config

import (
	"flag"
	"strings"
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.Func("b", "server address before short URL", func(s string) error {
		if !strings.HasSuffix(s, "/") {
				s += "/"
			}
			FlagServerAddr = s
			return nil
		})
	flag.Parse()
}
