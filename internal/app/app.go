package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sihce_diagnosticos/internal/config"
	"sihce_diagnosticos/internal/database"
	"sihce_diagnosticos/internal/modules"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Aplicacion encapsula todos los componentes de la aplicación.
type Aplicacion struct {
	servidor *http.Server
	db       *sql.DB
}

// Nueva crea y configura una nueva instancia de la aplicación.
func App() (*Aplicacion, error) {
	log.Println("Iniciando aplicación...")

	cfg := IniciarEnv()
	db, err := IniciarDatabase(cfg)
	if err != nil {
		return nil, err
	}

	enrutador := IniciarEnrutador(db)
	servidor := configurarServidor(cfg, enrutador)

	return &Aplicacion{
		servidor: servidor,
		db:       db,
	}, nil
}

// Ejecutar inicia el servidor HTTP y maneja el apagado gradual.
func (a *Aplicacion) Ejecutar() {
	defer database.CerrarConexion(a.db)

	go func() {
		log.Printf("✓ Servidor escuchando en el puerto %s", a.servidor.Addr)
		if err := a.servidor.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ No se pudo iniciar el servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando el servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.servidor.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Error en el apagado del servidor: %v", err)
	}

	log.Println("✓ Servidor apagado correctamente")
}

// --- Funciones auxiliares de configuración ---

func IniciarEnv() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, se usarán las variables de entorno del sistema")
	}
	cfg := config.CargarConfiguracion()
	return &cfg
}

func IniciarDatabase(cfg *config.Config) (*sql.DB, error) {
db, err := database.ConectarDB(*cfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func IniciarEnrutador(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Servir archivos estáticos para la ruta /public/
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// Registrar todos los módulos de la aplicación
	modules.IniciarModulos(db, r)

	return r
}

func configurarServidor(cfg *config.Config, handler http.Handler) *http.Server {
	port := cfg.ServerPort
	if port == "" {
		port = "8080" // Puerto por defecto
	}
	serverAddr := ":" + port

	return &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
}
