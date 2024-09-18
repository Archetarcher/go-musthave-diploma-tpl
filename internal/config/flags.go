package config

import (
	"flag"
)

const (
	flagRunAddrName                = "a"
	flagDatabaseUriName            = "d"
	flagAccrualSystemAddressName   = "r"
	flagTokenKeyName               = "k"
	flagTokenExpiresName           = "e"
	flagMigrationsPathName         = "m"
	flagWorkersCountName           = "w"
	flagWorkersRequestIntervalName = "i"
)

func (c *AppConfig) parseFlags() {
	flag.Parse()
}
func (c *AppConfig) initFlags() {

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.DatabaseUri, flagDatabaseUriName, "postgres://postgres:postgres@127.0.0.1:5432/gophermart?sslmode=disable", "database uri")
	flag.StringVar(&c.AccrualSystemAddress, flagAccrualSystemAddressName, "", "accrual system address")
	flag.StringVar(&c.Token.Key, flagTokenKeyName, "Zswx2zqELD", "key for token")
	flag.IntVar(&c.Token.ExpiresInMinutes, flagTokenExpiresName, 60, "token expires in minute")
	flag.StringVar(&c.MigrationsPath, flagMigrationsPathName, "internal/migrations", "migrations path")
	flag.IntVar(&c.Worker.Count, flagWorkersCountName, 4, "workers count")
	flag.IntVar(&c.Worker.Interval, flagWorkersRequestIntervalName, 2, "workers request interval")

}
