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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func CargarConfiguracion() Config {
	return Config{
		ServerPort:               getEnv("PORT", "3211"),
		DBUser:                   os.Getenv("DB_USER"),
		DBPassword:               os.Getenv("DB_PASSWORD"),
		DBServer:                 os.Getenv("DB_SERVER"),
		DBPort:                   getEnv("DB_PORT", "1433"),
		DBName:                   os.Getenv("DB_1"),
		DBEncrypt:                getEnv("DB_ENCRYPT", "false"),
		DBTrustServerCertificate: getEnv("DB_TRUST_SERVER_CERTIFICATE", "false"),
	}
}
