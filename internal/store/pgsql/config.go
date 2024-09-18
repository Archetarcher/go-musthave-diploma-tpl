package pgsql

type Config struct {
	Active         bool
	DatabaseURI    string
	MigrationsPath string
}
