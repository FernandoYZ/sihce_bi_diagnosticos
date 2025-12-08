package diagnostico

import (
	"context"
	"database/sql"
	"log"
	"sihce_diagnosticos/internal/models"
)

type repositorioDiagnostico interface {
	ObtenerDiagnosticos(ctx context.Context, pagina, cantidad int, buscar string) ([]models.Diagnostico, error)
}

type diagnosticoRepository struct {
	db *sql.DB
}

func DiagnosticoRepository(db *sql.DB) repositorioDiagnostico {
	return &diagnosticoRepository{db: db}
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
