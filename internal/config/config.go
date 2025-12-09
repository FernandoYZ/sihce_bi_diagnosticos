package config

import "os"

type Config struct {
	ServerPort               string
	DBUser                   string
	DBPassword               string
	DBServer                 string
	DBPort                   string
	DBName                   string
	DBEncrypt                string
	DBTrustServerCertificate string
}

func CargarConfiguracion() Config {
	return Config{
		ServerPort:               os.Getenv("PORT"),
		DBUser:                   os.Getenv("DB_USER"),
		DBPassword:               os.Getenv("DB_PASSWORD"),
		DBServer:                 os.Getenv("DB_SERVER"),
		DBPort:                   os.Getenv("DB_PORT"),
		DBName:                   os.Getenv("DB_1"),
		DBEncrypt:                os.Getenv("DB_ENCRYPT"),
		DBTrustServerCertificate: os.Getenv("DB_TRUST_SERVER_CERTIFICATE"),
	}
}
