package router

import (
	"net/http"
	"sihce_diagnosticos/internal/controller"
	"github.com/gorilla/mux"
)

func SetupRouter(diagnosticoController *controller.DiagnosticoController) *mux.Router {
	r := mux.NewRouter()

	// Servir archivos estáticos desde la carpeta 'public'
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// Rutas
	r.HandleFunc("/", diagnosticoController.PaginaHome).Methods(http.MethodGet)
	r.HandleFunc("/api/diagnosticos", diagnosticoController.ObtenerDiagnosticos).Methods(http.MethodGet)

	return r
}
