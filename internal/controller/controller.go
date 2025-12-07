package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sihce_diagnosticos/internal/repository"
	"sihce_diagnosticos/internal/views"
	"sihce_diagnosticos/internal/views/components"
	"strconv"
)

type DiagnosticoController struct {
	repo repository.DiagnosticoRepository
}

func NewDiagnosticoController(repo repository.DiagnosticoRepository) *DiagnosticoController {
	return &DiagnosticoController{repo: repo}
}

// ObtenerDiagnosticos maneja la solicitud GET para obtener los diagnósticos
func (c *DiagnosticoController) ObtenerDiagnosticos(w http.ResponseWriter, r *http.Request) {
	// Leer parámetros de la consulta
	paginaStr := r.URL.Query().Get("pagina")
	cantidadStr := r.URL.Query().Get("cantidad")
	buscar := r.URL.Query().Get("buscar")

	// Convertir los parámetros de página y cantidad
	pagina, err := strconv.Atoi(paginaStr)
	if err != nil || pagina < 1 {
		pagina = 1
	}
	cantidad, err := strconv.Atoi(cantidadStr)
	if err != nil || cantidad < 1 {
		cantidad = 10
	}

	// Obtener los diagnósticos desde el repositorio
	diagnosticos, err := c.repo.ObtenerDiagnosticos(r.Context(), pagina, cantidad, buscar)
	if err != nil {
		log.Printf("Error al obtener diagnósticos del repositorio: %v", err)
		http.Error(w, "Error al obtener los diagnósticos", http.StatusInternalServerError)
		return
	}

	// Si es una petición de HTMX, renderizar los componentes
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")

		// Renderizar la lista de diagnósticos
		startIndex := (pagina - 1) * cantidad
		if err := components.DiagnosticoList(diagnosticos, buscar, startIndex).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar DiagnosticoList: %v", err)
			http.Error(w, "Error al renderizar la lista", http.StatusInternalServerError)
			return
		}

		// Si obtuvimos una página completa, renderizar el botón "Cargar más"
		if len(diagnosticos) == cantidad {
			if err := components.LoadMoreDiagnosticos(pagina+1, buscar).Render(r.Context(), w); err != nil {
				log.Printf("Error al renderizar LoadMoreDiagnosticos: %v", err)
				http.Error(w, "Error al renderizar el botón de carga", http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("Fin de los resultados")
		}
		return
	}

	// Si no, devolver JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(diagnosticos); err != nil {
		http.Error(w, "Error al serializar los resultados", http.StatusInternalServerError)
	}
}

// RenderHomePage handles the request to render the home page
func (c *DiagnosticoController) PaginaHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Home("World").Render(r.Context(), w)
}

