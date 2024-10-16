package config

import (
	"flag"
)

const (
	flagRunAddrName                = "a"
	flagDatabaseURIName            = "d"
	flagAccrualSystemAddressName   = "r"
	flagTokenKeyName               = "k"
	flagTokenExpiresName           = "e"
	flagMigrationsPathName         = "m"
	flagWorkersCountName           = "w"
	flagWorkersRequestIntervalName = "i"
	flagRetryAfterName             = "t"
	flagPollIntervalName           = "p"
	flagRetryCountName             = "rc"
)

func (c *AppConfig) parseFlags() {
	flag.Parse()
}
func (c *AppConfig) initFlags() {

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.DatabaseURI, flagDatabaseURIName, "postgres://postgres:postgres@127.0.0.1:5432/gophermart?sslmode=disable", "database uri")
	flag.StringVar(&c.AccrualSystemAddress, flagAccrualSystemAddressName, "", "accrual system address")
	flag.StringVar(&c.Token.Key, flagTokenKeyName, "Zswx2zqELD", "key for token")
	flag.IntVar(&c.Token.ExpiresInMinutes, flagTokenExpiresName, 60, "token expires in minute")
	flag.StringVar(&c.MigrationsPath, flagMigrationsPathName, "internal/migrations", "migrations path")
	flag.IntVar(&c.Worker.Count, flagWorkersCountName, 4, "workers count")
	flag.IntVar(&c.Worker.Interval, flagWorkersRequestIntervalName, 2, "workers request interval")
	flag.IntVar(&c.RetryAfter, flagRetryAfterName, 60, "retry after")
	flag.IntVar(&c.PollInterval, flagPollIntervalName, 5, "interval in seconds for poll ")
	flag.IntVar(&c.RetryCount, flagRetryCountName, 3, "retry count")

}
