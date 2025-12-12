package diagnostico

import (
	"context"
	"database/sql"
	"log"
	"sihce_diagnosticos/internal/models"
)

type repositorioDiagnostico interface {
	ObtenerDiagnosticos(ctx context.Context, pagina, cantidad int, buscar string) ([]models.Diagnostico, error)
	GetResumenDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) (*models.ResumenDiagnostico, error)
	GetSexoPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.SexoPorDiagnostico, error)
	GetEdadesPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.EdadesPorDiagnostico, error)
	GetDistritosPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.DistritosPorDiagnostico, error)
	GetAtencionesPorTiempo(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin, query string) ([]models.AtencionesPorTiempo, error)
	ObtenerRangoEdadesSexo(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.RangoEdadSexo, error)
	ObtenerCondicionPaciente(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.CondicionPaciente, error)
	ObtenerClasificacionDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.ClasificacionDiagnostico, error)
}

type diagnosticoRepository struct {
	db *sql.DB
}

func DiagnosticoRepository(db *sql.DB) repositorioDiagnostico {
	return &diagnosticoRepository{db: db}
}

func (r *diagnosticoRepository) GetResumenDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) (*models.ResumenDiagnostico, error) {
	var resumen models.ResumenDiagnostico
	err := r.db.QueryRowContext(ctx, QUERY_RESUMEN_CARDS,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	).Scan(
		&resumen.DistritosAfectadosMesActual,
		&resumen.DistritosAfectadosMesAnterior,
		&resumen.DiferenciaDistritosAfectados,
		&resumen.TotalPacientesUnicosActual,
		&resumen.TotalPacientesUnicosMesAnterior,
		&resumen.PorcentajeCambioPacientes,
		&resumen.TotalAtencionesMesActual,
		&resumen.TotalAtencionesMesAnterior,
		&resumen.PorcentajeCambioAtenciones,
		&resumen.RatioDeRetorno,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No results found is not an error here
		}
		log.Printf("Error executing or scanning summary query: %v", err)
		return nil, err
	}
	return &resumen, nil
}

func (r *diagnosticoRepository) GetSexoPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.SexoPorDiagnostico, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_SEXO_POR_DIAGNOSTICO,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.SexoPorDiagnostico
	for rows.Next() {
		var res models.SexoPorDiagnostico
		if err := rows.Scan(&res.Sexo, &res.CantidadAtenciones); err != nil {
			log.Printf("Error al escanear resultado de SexoPorDiagnostico: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de SexoPorDiagnostico: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) GetEdadesPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.EdadesPorDiagnostico, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_EDADES_POR_DIAGNOSTICO,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.EdadesPorDiagnostico
	for rows.Next() {
		var res models.EdadesPorDiagnostico
		if err := rows.Scan(&res.GrupoEdad, &res.CantidadAtenciones); err != nil {
			log.Printf("Error al escanear resultado de EdadesPorDiagnostico: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de EdadesPorDiagnostico: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) GetDistritosPorDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.DistritosPorDiagnostico, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_CANTIDADES_POR_DISTRITO,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.DistritosPorDiagnostico
	for rows.Next() {
		var res models.DistritosPorDiagnostico
		if err := rows.Scan(&res.IdDistrito, &res.NombreDistrito, &res.NombreProvincia, &res.CantidadAtenciones); err != nil {
			log.Printf("Error al escanear resultado de DistritosPorDiagnostico: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de DistritosPorDiagnostico: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) GetAtencionesPorTiempo(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin, query string) ([]models.AtencionesPorTiempo, error) {
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.AtencionesPorTiempo
	for rows.Next() {
		var res models.AtencionesPorTiempo
		if err := rows.Scan(&res.Fecha, &res.CantidadAtenciones); err != nil {
			log.Printf("Error al escanear resultado de AtencionesPorTiempo: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de AtencionesPorTiempo: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) ObtenerDiagnosticos(ctx context.Context, pagina, cantidad int, buscar string) ([]models.Diagnostico, error) {
	offset := (pagina - 1) * cantidad

	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_DIAGNOSTICOS,
		sql.Named("offset", offset),
		sql.Named("cantidad", cantidad),
		sql.Named("buscar", buscar),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diagnosticos []models.Diagnostico
	for rows.Next() {
		var diagnostico models.Diagnostico
		if err := rows.Scan(&diagnostico.IdDiagnostico, &diagnostico.Diagnostico); err != nil {
			log.Printf("Error al escanear resultado: %v", err)
			continue
		}
		diagnosticos = append(diagnosticos, diagnostico)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración: %v", err)
		return nil, err
	}

	return diagnosticos, nil
}

func (r *diagnosticoRepository) ObtenerRangoEdadesSexo(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.RangoEdadSexo, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_RANGO_EDADES_SEXO,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.RangoEdadSexo
	for rows.Next() {
		var res models.RangoEdadSexo
		if err := rows.Scan(&res.RangoEdad, &res.Masculino, &res.Femenino); err != nil {
			log.Printf("Error al escanear resultado de RangoEdadSexo: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de RangoEdadSexo: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) ObtenerCondicionPaciente(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.CondicionPaciente, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_CONDICION_PACIENTE,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.CondicionPaciente
	for rows.Next() {
		var res models.CondicionPaciente
		if err := rows.Scan(&res.TipoCondicionAlServicio, &res.Cantidad); err != nil {
			log.Printf("Error al escanear resultado de CondicionPaciente: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de CondicionPaciente: %v", err)
		return nil, err
	}
	return resultados, nil
}

func (r *diagnosticoRepository) ObtenerClasificacionDiagnostico(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.ClasificacionDiagnostico, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_CONDICION_PACIENTE,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.ClasificacionDiagnostico
	for rows.Next() {
		var res models.ClasificacionDiagnostico
		if err := rows.Scan(&res.ClasificacionDiagnostico, &res.Cantidad); err != nil {
			log.Printf("Error al escanear resultado de ClasificacionDiagnostico: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de ClasificacionDiagnostico: %v", err)
		return nil, err
	}
	return resultados, nil
}