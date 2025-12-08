package diagnostico

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// RegisterModule inicializa y registra todos los componentes del módulo de diagnóstico.
func DiagnosticoModule(db *sql.DB, router *mux.Router) {
	// 1. Inicializar Repositorio
	repo := DiagnosticoRepository(db)

	// 2. Inicializar Servicio
	servicio := DiagnosticoService(repo)

	// 3. Inicializar Controlador
	controller := DiagnosticoController(servicio)

	// 4. Registrar las rutas del módulo
	controller.Router(router)
}
