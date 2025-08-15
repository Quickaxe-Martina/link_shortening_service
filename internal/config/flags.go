package config

import (
	"flag"
)

func ParseFlags() {

	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&FlagServerAddr, "b", "http://localhost:8080/", "server address before short URL")
	flag.Parse()
}
