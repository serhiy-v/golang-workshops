package config

type Application struct {
	DB *Database
}

type Database struct {
	DSN string `env:"DSN"`
}
