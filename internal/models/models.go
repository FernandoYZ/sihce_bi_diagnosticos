package models

import "database/sql"

type Diagnostico struct {
	IdDiagnostico int    `json:"idDiagnostico"`
	Diagnostico   string `json:"diagnostico"`
}

type ResumenDiagnostico struct {
	DistritosAfectadosMesActual     sql.NullInt64   `json:"distritosAfectadosMesActual"`
	DistritosAfectadosMesAnterior   sql.NullInt64   `json:"distritosAfectadosMesAnterior"`
	DiferenciaDistritosAfectados    sql.NullInt64   `json:"diferenciaDistritosAfectados"`
	TotalPacientesUnicosActual      sql.NullInt64   `json:"totalPacientesUnicosActual"`
	TotalPacientesUnicosMesAnterior sql.NullInt64   `json:"totalPacientesUnicosMesAnterior"`
	PorcentajeCambioPacientes       sql.NullFloat64 `json:"porcentajeCambioPacientes"`
	TotalAtencionesMesActual        sql.NullInt64   `json:"totalAtencionesMesActual"`
	TotalAtencionesMesAnterior      sql.NullInt64   `json:"totalAtencionesMesAnterior"`
	PorcentajeCambioAtenciones      sql.NullFloat64 `json:"porcentajeCambioAtenciones"`
	RatioDeRetorno                  sql.NullFloat64 `json:"ratioDeRetorno"`
}

type SexoPorDiagnostico struct {
	Sexo               string `json:"sexo"`
	CantidadAtenciones int    `json:"cantidadAtenciones"`
}

type EdadesPorDiagnostico struct {
	GrupoEdad          string `json:"grupoEdad"`
	CantidadAtenciones int    `json:"cantidadAtenciones"`
}

type DistritosPorDiagnostico struct {
	IdDistrito         int    `json:"idDistrito"`
	NombreDistrito     string `json:"nombreDistrito"`
	NombreProvincia    string `json:"nombreProvincia"`
	CantidadAtenciones int    `json:"cantidadAtenciones"`
}

type AtencionesPorTiempo struct {
	Fecha              string `json:"fecha"`
	CantidadAtenciones int    `json:"cantidadAtenciones"`
}

type AtencionesTiempoResponse struct {
	PeriodoActual   []AtencionesPorTiempo `json:"periodoActual"`
	PeriodoAnterior []AtencionesPorTiempo `json:"periodoAnterior"`
	Tipo            string                `json:"tipo"`
}

type RangoEdadSexo struct {
	RangoEdad string `json:"rangoEdad"`
	Masculino string `json:"masculino"`
	Femenino  string `json:"femenino"`
}

type CondicionPaciente struct {
	TipoCondicionAlServicio string `json:"condicionPaciente"`
	Cantidad      int    `json:"cantidad"`
}

type ClasificacionDiagnostico struct {
	ClasificacionDiagnostico string `json:"clasificacionDiagnostico"`
	Cantidad      int    `json:"cantidad"`
}
