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
	GetAtencionesPorDia(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.AtencionesPorDia, error)
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

func (r *diagnosticoRepository) GetAtencionesPorDia(ctx context.Context, idDiagnostico int, fechaInicio, fechaFin string) ([]models.AtencionesPorDia, error) {
	rows, err := r.db.QueryContext(ctx, QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_DIA,
		sql.Named("IdDiagnostico", idDiagnostico),
		sql.Named("FechaInicio", fechaInicio),
		sql.Named("FechaFin", fechaFin),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []models.AtencionesPorDia
	for rows.Next() {
		var res models.AtencionesPorDia
		if err := rows.Scan(&res.Fecha, &res.CantidadAtenciones); err != nil {
			log.Printf("Error al escanear resultado de AtencionesPorDia: %v", err)
			continue
		}
		resultados = append(resultados, res)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error durante la iteración de AtencionesPorDia: %v", err)
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
