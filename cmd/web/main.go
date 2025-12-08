package main

import (
	"log"
	"sihce_diagnosticos/internal/app"
)

func main() {
	aplicacion, err := app.App()
	if err != nil {
		log.Fatalf("❌ Error al inicializar la aplicación: %v", err)
	}

	aplicacion.Ejecutar()
}