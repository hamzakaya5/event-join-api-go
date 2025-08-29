package models

type PostgreConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	MaxConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
