package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig
	Rundeck RundeckConfig
	Hooks   []HookConfig
}

type ServerConfig struct {
	Host string
	Port uint16
}

type RundeckConfig struct {
	URL       string
	AuthToken string `toml:"auth_token"`
}

type HookConfig struct {
	URL    string
	Branch string
	JobID  string `toml:"job_id"`
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	config := Config{
		Server: ServerConfig{
			Host: "",
			Port: 8080,
		},
	}
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		logger.Panicln(err)
	}
	logger.Printf("Config: %+v\n", config)

	server, err := NewServer(&config, logger)
	if err != nil {
		logger.Panicln(err)
	}
	server.Start()
}
