package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sihce_diagnosticos/internal/config"
	"sihce_diagnosticos/internal/controller"
	"sihce_diagnosticos/internal/database"
	"sihce_diagnosticos/internal/repository"
	"sihce_diagnosticos/internal/router"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Iniciando servidor...")

	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, se usarán las variables de entorno del sistema")
	}

	log.Printf("DB_TRUST_SERVER_CERTIFICATE from config: %s", os.Getenv("DB_TRUST_SERVER_CERTIFICATE"))

	// 1. Configuración
	cfg := config.CargarConfiguracion()

	// 2. Conexión a la base de datos
	db, err := database.ConectarDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer database.CerrarConexion(db)

	// 3. Repositorio
	diagnosticoRepo := repository.NewDiagnosticoRepository(db)

	// 4. Controlador
	diagnosticoController := controller.NewDiagnosticoController(diagnosticoRepo)

	// 5. Router
	r := router.SetupRouter(diagnosticoController)

	// 6. Servidor HTTP
	port := cfg.ServerPort
	if port == "" {
		port = "8080" // Puerto por defecto
	}
	serverAddr := ":" + port
	
	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		log.Printf("✓ Servidor escuchando en el puerto %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ No se pudo iniciar el servidor: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando el servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Error en el apagado del servidor: %v", err)
	}

	log.Println("✓ Servidor apagado correctamente")
}