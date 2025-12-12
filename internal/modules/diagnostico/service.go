package diagnostico

import (
	"context"
	"sihce_diagnosticos/internal/models"
	"time"
)

// Estructura de resultado que retornará la función
type parametrosValidos struct {
	idDiagnostico  int
	fechaInicio    time.Time
	fechaFin       time.Time
	fechaInicioFmt string
	fechaFinFmt    string
}

type servicioDiagnostico struct {
	repo repositorioDiagnostico
}

func DiagnosticoService(repo repositorioDiagnostico) *servicioDiagnostico {
	return &servicioDiagnostico{repo: repo}
}

// Función para validar y preparar los parámetros (ID y fechas)
func (s *servicioDiagnostico) prepararParámetros(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) (*parametrosValidos, error) {
	// Validar los parámetros
	if err := ValidarParametrosDiagnostico(idDiagnosticoStr, fechaInicioStr, fechaFinStr); err != nil {
		return nil, err
	}

	// Convertir ID
	idDiagnostico, err := ConvertirIdDiagnostico(idDiagnosticoStr)
	if err != nil {
		return nil, err
	}

	// Parsear las fechas
	fechaInicio, fechaFin, err := ParsearFechas(fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	// Formatear las fechas para SQL
	fechaInicioFmt := FormatearFechaSQL(fechaInicio)
	fechaFinFmt := FormatearFechaSQL(fechaFin)

	// Retornar el resultado
	return &parametrosValidos{
		idDiagnostico:  idDiagnostico,
		fechaInicio:    fechaInicio,
		fechaFin:       fechaFin,
		fechaInicioFmt: fechaInicioFmt,
		fechaFinFmt:    fechaFinFmt,
	}, nil
}

func (s *servicioDiagnostico) ObtenerDiagnosticos(ctx context.Context, pagina int, cantidad int, buscar string) ([]models.Diagnostico, error) {
	return s.repo.ObtenerDiagnosticos(ctx, pagina, cantidad, buscar)
}

// GetResumenDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetResumenDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) (*models.ResumenDiagnostico, error) {
	// Llamar a la función común para validación y preparación de parámetros
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	// Llamar al repositorio con los parámetros ya preparados
	return s.repo.GetResumenDiagnostico(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

// GetSexoPorDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetSexoPorDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.SexoPorDiagnostico, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSexoPorDiagnostico(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

// GetEdadesPorDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetEdadesPorDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.EdadesPorDiagnostico, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetEdadesPorDiagnostico(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

// GetDistritosPorDiagnosticoConValidacion valida, parsea y formatea las fechas antes de llamar al repositorio
func (s *servicioDiagnostico) GetDistritosPorDiagnosticoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.DistritosPorDiagnostico, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDistritosPorDiagnostico(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

var mesesCortos = []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}

func formatarFechaParaGrafico(fecha time.Time, tipo string) string {
	switch tipo {
	case "YoY":
		return fecha.Format("2006")
	case "MoM":
		return mesesCortos[fecha.Month()-1]
	case "WoW", "DoD":
		return fecha.Format("02") + " " + mesesCortos[fecha.Month()-1]
	default:
		return fecha.Format("2006-01-02")
	}
}

func (s *servicioDiagnostico) GetAtencionesPorDiaConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) (*models.AtencionesTiempoResponse, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}

	diferenciaDias := int(parametros.fechaFin.Sub(parametros.fechaInicio).Hours() / 24)
	var query, tipo string

	switch {
	case diferenciaDias <= 31:
		query = QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_DIA
		tipo = "DoD"
	case diferenciaDias <= 120:
		query = QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_SEMANA
		tipo = "WoW"
	case diferenciaDias <= 730:
		query = QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_MES
		tipo = "MoM"
	default:
		query = QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_ANIO
		tipo = "YoY"
	}

	// Periodo actual
	periodoActualData, err := s.repo.GetAtencionesPorTiempo(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt, query)
	if err != nil {
		return nil, err
	}

	// Periodo anterior
	diferencia := parametros.fechaFin.Sub(parametros.fechaInicio)
	fechaInicioAnterior := parametros.fechaInicio.Add(-diferencia).Add(-24 * time.Hour)
	fechaFinAnterior := parametros.fechaInicio.Add(-24 * time.Hour)
	
	periodoAnteriorData, err := s.repo.GetAtencionesPorTiempo(ctx, parametros.idDiagnostico, FormatearFechaSQL(fechaInicioAnterior), FormatearFechaSQL(fechaFinAnterior), query)
	if err != nil {
		return nil, err
	}

	datosActualesMap := make(map[string]int)
	for _, v := range periodoActualData {
		fecha, _ := time.Parse("2006-01-02T15:04:05Z", v.Fecha)
		datosActualesMap[fecha.Format("2006-01-02")] = v.CantidadAtenciones
	}
	datosAnterioresMap := make(map[string]int)
	for _, v := range periodoAnteriorData {
		fecha, _ := time.Parse("2006-01-02T15:04:05Z", v.Fecha)
		datosAnterioresMap[fecha.Format("2006-01-02")] = v.CantidadAtenciones
	}

	var periodoActualProcesado []models.AtencionesPorTiempo
	var periodoAnteriorProcesado []models.AtencionesPorTiempo

	current := parametros.fechaInicio
	for current.Before(parametros.fechaFin) || current.Equal(parametros.fechaFin) {
		fechaStr := current.Format("2006-01-02")
		cantidadActual, ok := datosActualesMap[fechaStr]
		if !ok {
			cantidadActual = 0
		}
		periodoActualProcesado = append(periodoActualProcesado, models.AtencionesPorTiempo{
			Fecha:              formatarFechaParaGrafico(current, tipo),
			CantidadAtenciones: cantidadActual,
		})

		fechaAnteriorEquivalente := current.Add(-diferencia).Add(-24 * time.Hour)
		fechaAnteriorStr := fechaAnteriorEquivalente.Format("2006-01-02")
		cantidadAnterior, ok := datosAnterioresMap[fechaAnteriorStr]
		if !ok {
			cantidadAnterior = 0
		}
		periodoAnteriorProcesado = append(periodoAnteriorProcesado, models.AtencionesPorTiempo{
			Fecha:              formatarFechaParaGrafico(current, tipo),
			CantidadAtenciones: cantidadAnterior,
		})

		switch tipo {
		case "DoD":
			current = current.AddDate(0, 0, 1)
		case "WoW":
			current = current.AddDate(0, 0, 7)
		case "MoM":
			current = current.AddDate(0, 1, 0)
		case "YoY":
			current = current.AddDate(1, 0, 0)
		default:
			current = current.AddDate(0, 0, 1)
		}
	}

	return &models.AtencionesTiempoResponse{
		PeriodoActual:   periodoActualProcesado,
		PeriodoAnterior: periodoAnteriorProcesado,
		Tipo:            tipo,
	}, nil
}


func (s *servicioDiagnostico) GetRangoEdadesSexoConValidacion(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.RangoEdadSexo, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ObtenerRangoEdadesSexo(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

func (s *servicioDiagnostico) ObtenerCondicionPaciente(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.CondicionPaciente, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ObtenerCondicionPaciente(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}

func (s *servicioDiagnostico) ObtenerClasificacionDiagnostico(ctx context.Context, idDiagnosticoStr, fechaInicioStr, fechaFinStr string) ([]models.ClasificacionDiagnostico, error) {
	parametros, err := s.prepararParámetros(ctx, idDiagnosticoStr, fechaInicioStr, fechaFinStr)
	if err != nil {
		return nil, err
	}
	return s.repo.ObtenerClasificacionDiagnostico(ctx, parametros.idDiagnostico, parametros.fechaInicioFmt, parametros.fechaFinFmt)
}