package config

import "os"

type ConfigBD struct {
	Usuario					string
	Contrasena				string
	Host					string
	Puerto					string
	NombreBD				string
	Encrypt					string
	TrustServerCertificate	string
}

func ConfiguracionBD() ConfigBD {
	usuario := os.Getenv("DB_USER")
	contrasena := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_SERVER")
	puerto := os.Getenv("DB_PORT")
	nombreBD := os.Getenv("DB_1")
	encrypt := os.Getenv("DB_ENCRYPT")
	trustServerCertificate := os.Getenv("DB_TRUST_SERVER_CERTIFICATE")

	return ConfigBD{
		Usuario:    usuario,
		Contrasena: contrasena,
		Host:       host,
		Puerto:     puerto,
		NombreBD:   nombreBD,
		Encrypt:    encrypt,
		TrustServerCertificate: trustServerCertificate,
	}
}