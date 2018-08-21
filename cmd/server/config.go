package main

import (
	"github.com/jinzhu/configor"
)

// Application config structure
type config struct {
	DBSource      string
	ListenAddress string
}

// Loading the app configuration
func loadConfig() *config {
	config := &config{}
	if err := configor.Load(config, "/config/app.yml"); err != nil {
		panic(err)
	}
	return config
}
