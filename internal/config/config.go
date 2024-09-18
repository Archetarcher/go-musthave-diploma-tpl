package config

import "github.com/go-chi/jwtauth/v5"

func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

type AppConfig struct {
	RunAddr              string
	DatabaseUri          string
	MigrationsPath       string
	AccrualSystemAddress string
	Token                Token
	Worker               Worker
}
type Token struct {
	Key              string
	ExpiresInMinutes int
	AuthToken        *jwtauth.JWTAuth
}
type Worker struct {
	Count    int
	Interval int
}

func NewConfig() *AppConfig {
	var c = AppConfig{}
	c.initFlags()

	return &c
}
