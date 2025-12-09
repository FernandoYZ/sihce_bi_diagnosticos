package diagnostico

import (
	"encoding/json"
	"log"
	"net/http"
	"sihce_diagnosticos/internal/views"
	"sihce_diagnosticos/internal/views/components"
	"strconv"

	"github.com/gorilla/mux"
)

type controller struct {
	servicio *servicioDiagnostico
}

func DiagnosticoController(servicio *servicioDiagnostico) *controller {
	return &controller{servicio: servicio}
}

func (c *controller) Router(r *mux.Router) {
	r.HandleFunc("/", c.PaginaHome).Methods(http.MethodGet)
	r.HandleFunc("/api/diagnosticos", c.ObtenerDiagnosticos).Methods(http.MethodGet)
	r.HandleFunc("/api/resumen", c.GetResumenHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/sexo-por-diagnostico", c.GetSexoPorDiagnosticoHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/edades-por-diagnostico", c.GetEdadesPorDiagnosticoHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/distritos-por-diagnostico", c.GetDistritosPorDiagnosticoHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/atenciones-por-dia", c.GetAtencionesPorDiaHandler).Methods(http.MethodGet)
}

func (c *controller) GetSexoPorDiagnosticoHandler(w http.ResponseWriter, r *http.Request) {
	idDiagnosticoStr := r.URL.Query().Get("IdDiagnostico")
	fechaInicioStr := r.URL.Query().Get("FechaInicio")
	fechaFinStr := r.URL.Query().Get("FechaFin")

	sexoData, err := c.servicio.GetSexoPorDiagnosticoConValidacion(r.Context(), idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		log.Printf("Error en GetSexoPorDiagnosticoHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sexoData); err != nil {
		http.Error(w, "Error al serializar los resultados de sexo por diagnostico", http.StatusInternalServerError)
	}
}

func (c *controller) GetEdadesPorDiagnosticoHandler(w http.ResponseWriter, r *http.Request) {
	idDiagnosticoStr := r.URL.Query().Get("IdDiagnostico")
	fechaInicioStr := r.URL.Query().Get("FechaInicio")
	fechaFinStr := r.URL.Query().Get("FechaFin")

	edadesData, err := c.servicio.GetEdadesPorDiagnosticoConValidacion(r.Context(), idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		log.Printf("Error en GetEdadesPorDiagnosticoHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(edadesData); err != nil {
		http.Error(w, "Error al serializar los resultados de edades por diagnostico", http.StatusInternalServerError)
	}
}

func (c *controller) GetDistritosPorDiagnosticoHandler(w http.ResponseWriter, r *http.Request) {
	idDiagnosticoStr := r.URL.Query().Get("IdDiagnostico")
	fechaInicioStr := r.URL.Query().Get("FechaInicio")
	fechaFinStr := r.URL.Query().Get("FechaFin")

	distritosData, err := c.servicio.GetDistritosPorDiagnosticoConValidacion(r.Context(), idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		log.Printf("Error en GetDistritosPorDiagnosticoHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(distritosData); err != nil {
		http.Error(w, "Error al serializar los resultados de distritos por diagnostico", http.StatusInternalServerError)
	}
}

func (c *controller) GetAtencionesPorDiaHandler(w http.ResponseWriter, r *http.Request) {
	idDiagnosticoStr := r.URL.Query().Get("IdDiagnostico")
	fechaInicioStr := r.URL.Query().Get("FechaInicio")
	fechaFinStr := r.URL.Query().Get("FechaFin")

	atencionesPorDiaData, err := c.servicio.GetAtencionesPorDiaConValidacion(r.Context(), idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		log.Printf("Error en GetAtencionesPorDiaHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(atencionesPorDiaData); err != nil {
		http.Error(w, "Error al serializar los resultados de atenciones por dia", http.StatusInternalServerError)
	}
}

func (c *controller) GetResumenHandler(w http.ResponseWriter, r *http.Request) {
	idDiagnosticoStr := r.URL.Query().Get("IdDiagnostico")
	fechaInicioStr := r.URL.Query().Get("FechaInicio")
	fechaFinStr := r.URL.Query().Get("FechaFin")

	resumen, err := c.servicio.GetResumenDiagnosticoConValidacion(r.Context(), idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		log.Printf("Error en GetResumenHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resumen == nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p>No se encontraron datos para el diagnóstico y rango de fechas seleccionado.</p>"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := components.SummaryCards(*resumen).Render(r.Context(), w); err != nil {
		log.Printf("Error al renderizar SummaryCards: %v", err)
		http.Error(w, "Error al renderizar el resumen", http.StatusInternalServerError)
		return
	}
}

// ObtenerDiagnosticos maneja la solicitud GET para obtener los diagnósticos
func (c *controller) ObtenerDiagnosticos(w http.ResponseWriter, r *http.Request) {
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

	// Obtener los diagnósticos desde el servicio
	diagnosticos, err := c.servicio.ObtenerDiagnosticos(r.Context(), pagina, cantidad, buscar)
	if err != nil {
		log.Printf("Error al obtener diagnósticos del servicio: %v", err)
		http.Error(w, "Error al obtener los diagnósticos", http.StatusInternalServerError)
		return
	}

	// Si es una petición de HTMX, renderizar los componentes
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		startIndex := (pagina - 1) * cantidad
		hasMore := len(diagnosticos) == cantidad

		// For infinite scroll requests (page > 1), use the OOB component
		if pagina > 1 {
			if err := components.DiagnosticoListOOB(diagnosticos, buscar, startIndex, pagina+1, hasMore).Render(r.Context(), w); err != nil {
				log.Printf("Error al renderizar DiagnosticoListOOB: %v", err)
				http.Error(w, "Error al renderizar la lista OOB", http.StatusInternalServerError)
			}
			return
		}

		// For initial load (page 1), render the list and the trigger normally
		if err := components.DiagnosticoList(diagnosticos, buscar, startIndex).Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar DiagnosticoList: %v", err)
			http.Error(w, "Error al renderizar la lista", http.StatusInternalServerError)
			return
		}

		if hasMore {
			if err := components.LoadMoreDiagnosticos(pagina+1, buscar).Render(r.Context(), w); err != nil {
				log.Printf("Error al renderizar LoadMoreDiagnosticos: %v", err)
				http.Error(w, "Error al renderizar el botón de carga", http.StatusInternalServerError)
				return
			}
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
func (c *controller) PaginaHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Home("World").Render(r.Context(), w)
}
