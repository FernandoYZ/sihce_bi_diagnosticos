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