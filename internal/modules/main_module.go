package modules

import (
	"database/sql"
	"sihce_diagnosticos/internal/modules/diagnostico"

	"github.com/gorilla/mux"
)

// RegisterAllModules llama a los registradores de todos los módulos de la aplicación.
func IniciarModulos(db *sql.DB, router *mux.Router) {
	// Registrar módulo de Diagnóstico
	diagnostico.DiagnosticoModule(db, router)

}
