package pgsql

type Config struct {
	Active         bool
	DatabaseUri    string
	MigrationsPath string
}
