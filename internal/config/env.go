package config

import (
	"os"
	"strconv"
)

const (
	envRunAddrName              = "RUN_ADDRESS"
	envDatabaseUriName          = "DATABASE_URI"
	envAccrualSystemAddressName = "ACCRUAL_SYSTEM_ADDRESS"
	envTokenKeyName             = "TOKEN_KEY"
	envTokenExpiresName         = "TOKEN_EXPIRES_IN_MINUTES"
	envWorkersCountName         = "WORKERS_COUNT"
	envWorkersIntervalName      = "WORKERS_INTERVAL"
)

func getEnvOrDefault(env string, def any, t int) any {
	val := os.Getenv(env)
	if val == "" {
		return def
	}

	switch t {
	case 1:
		return val
	case 2:
		if i, err := strconv.Atoi(val); err == nil {
			return i

		}
		return def
	case 3:
		if i, err := strconv.ParseBool(val); err == nil {
			return i
		}
		return def
	default:
		return def
	}
}

func (c *AppConfig) parseEnv() {
	c.RunAddr = getEnvOrDefault(envRunAddrName, c.RunAddr, 1).(string)
	c.DatabaseUri = getEnvOrDefault(envDatabaseUriName, c.DatabaseUri, 1).(string)
	c.AccrualSystemAddress = getEnvOrDefault(envAccrualSystemAddressName, c.AccrualSystemAddress, 1).(string)
	c.Token.Key = getEnvOrDefault(envTokenKeyName, c.Token.Key, 1).(string)
	c.Token.ExpiresInMinutes = getEnvOrDefault(envTokenExpiresName, c.Token.ExpiresInMinutes, 2).(int)
	c.Worker.Count = getEnvOrDefault(envWorkersCountName, c.Worker.Count, 2).(int)
	c.Worker.Interval = getEnvOrDefault(envWorkersIntervalName, c.Worker.Interval, 2).(int)

}
