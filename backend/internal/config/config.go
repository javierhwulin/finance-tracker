package config

import "flag"

type Config struct {
	Version string
	Port    string
	Env     string
}

func NewConfig() *Config {
	version := flag.String("version", "1.0.0", "Version")
	port := flag.String("port", "8080", "Port to listen on")
	env := flag.String("env", "development", "Environment")
	flag.Parse()
	return &Config{
		Version: *version,
		Port:    *port,
		Env:     *env,
	}
}
