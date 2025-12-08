package diagnostico

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrCamposRequeridos    = errors.New("todos los campos son requeridos: IdDiagnostico, FechaInicio, FechaFin")
	ErrIdDiagnosticoInvalido = errors.New("IdDiagnostico debe ser un número válido")
	ErrFechaInicioInvalida  = errors.New("formato de FechaInicio inválido")
	ErrFechaFinInvalida     = errors.New("formato de FechaFin inválido")
)

// ValidarParametrosDiagnostico valida los parámetros comunes de diagnóstico
func ValidarParametrosDiagnostico(idDiagnosticoStr, fechaInicioStr, fechaFinStr string) error {
	if idDiagnosticoStr == "" || fechaInicioStr == "" || fechaFinStr == "" {
		return ErrCamposRequeridos
	}
	return nil
}

// ConvertirIdDiagnostico convierte el string a entero
func ConvertirIdDiagnostico(idDiagnosticoStr string) (int, error) {
	idDiagnostico, err := strconv.Atoi(idDiagnosticoStr)
	if err != nil {
		return 0, ErrIdDiagnosticoInvalido
	}
	return idDiagnostico, nil
}

// ParsearFechas parsea y valida las fechas de entrada
func ParsearFechas(fechaInicioStr, fechaFinStr string) (time.Time, time.Time, error) {
	layoutOriginal := "2006-01-02"

	fechaInicio, err := time.Parse(layoutOriginal, fechaInicioStr)
	if err != nil {
		return time.Time{}, time.Time{}, ErrFechaInicioInvalida
	}

	fechaFin, err := time.Parse(layoutOriginal, fechaFinStr)
	if err != nil {
		return time.Time{}, time.Time{}, ErrFechaFinInvalida
	}

	return fechaInicio, fechaFin, nil
}

// FormatearFechaSQL convierte una fecha a formato YYYYMMDD para SQL Server
func FormatearFechaSQL(fecha time.Time) string {
	return fecha.Format("20060102")
}
