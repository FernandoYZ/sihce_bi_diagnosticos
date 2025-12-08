package diagnostico

import (
	"context"
	"sihce_diagnosticos/internal/models"
)

type servicioDiagnostico struct {
	repo repositorioDiagnostico
}

func DiagnosticoService(repo repositorioDiagnostico) *servicioDiagnostico {
	return &servicioDiagnostico{repo: repo}
}

func (s *servicioDiagnostico) ObtenerDiagnosticos(ctx context.Context, pagina int, cantidad int, buscar string) ([]models.Diagnostico, error) {
	return s.repo.ObtenerDiagnosticos(ctx, pagina, cantidad, buscar)
}

// GetResumenDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetResumenDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) (*models.ResumenDiagnostico, error) {
	// Validar parámetros
	if err := ValidarParametrosDiagnostico(idDiagnosticoStr, fechaInicioStr, fechaFinStr); err != nil {
		return nil, err
	}

	// Convertir ID
	idDiagnostico, err := ConvertirIdDiagnostico(idDiagnosticoStr)
	if err != nil {
		return nil, err
	}

	// Parsear y validar fechas
	fechaInicio, fechaFin, err := ParsearFechas(fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	// Formatear fechas para SQL
	fechaInicioFmt := FormatearFechaSQL(fechaInicio)
	fechaFinFmt := FormatearFechaSQL(fechaFin)

	return s.repo.GetResumenDiagnostico(ctx, idDiagnostico, fechaInicioFmt, fechaFinFmt)
}

// GetSexoPorDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetSexoPorDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.SexoPorDiagnostico, error) {
	// Validar parámetros
	if err := ValidarParametrosDiagnostico(idDiagnosticoStr, fechaInicioStr, fechaFinStr); err != nil {
		return nil, err
	}

	// Convertir ID
	idDiagnostico, err := ConvertirIdDiagnostico(idDiagnosticoStr)
	if err != nil {
		return nil, err
	}

	// Parsear y validar fechas
	fechaInicio, fechaFin, err := ParsearFechas(fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	// Formatear fechas para SQL
	fechaInicioFmt := FormatearFechaSQL(fechaInicio)
	fechaFinFmt := FormatearFechaSQL(fechaFin)

	return s.repo.GetSexoPorDiagnostico(ctx, idDiagnostico, fechaInicioFmt, fechaFinFmt)
}

// GetEdadesPorDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetEdadesPorDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.EdadesPorDiagnostico, error) {
	// Validar parámetros
	if err := ValidarParametrosDiagnostico(idDiagnosticoStr, fechaInicioStr, fechaFinStr); err != nil {
		return nil, err
	}

	// Convertir ID
	idDiagnostico, err := ConvertirIdDiagnostico(idDiagnosticoStr)
	if err != nil {
		return nil, err
	}

	// Parsear y validar fechas
	fechaInicio, fechaFin, err := ParsearFechas(fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	// Formatear fechas para SQL
	fechaInicioFmt := FormatearFechaSQL(fechaInicio)
	fechaFinFmt := FormatearFechaSQL(fechaFin)

	return s.repo.GetEdadesPorDiagnostico(ctx, idDiagnostico, fechaInicioFmt, fechaFinFmt)
}
