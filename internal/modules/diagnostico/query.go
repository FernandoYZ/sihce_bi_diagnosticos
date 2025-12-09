package diagnostico

const (
	QUERY_OBTENER_DIAGNOSTICOS = `
		SELECT 
			IdDiagnostico,
			RTRIM(CodigoCIE2004) + ' - ' + RTRIM(Descripcion) AS Diagnostico
		FROM Diagnosticos
		WHERE 
			CodigoCIE2004 IS NOT NULL
			AND Descripcion IS NOT NULL
			AND (
				@buscar = '' 
				OR (CodigoCIE2004 + ' - ' + Descripcion) LIKE '%' + @buscar + '%'
			)
		ORDER BY CodigoCIE2004
		OFFSET @offset ROWS
		FETCH NEXT @cantidad ROWS ONLY;
	`

	QUERY_RESUMEN_CARDS = `
		WITH CTE_AtencionesFiltradas AS (
			-- CTE para obtener las atenciones filtradas por diagnóstico, fecha y otros criterios
			SELECT 
				a.IdAtencion, 
				p.IdPaciente, 
				p.IdDistritoProcedencia, 
				a.FechaIngreso
			FROM Atenciones a
			INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
			INNER JOIN Diagnosticos d ON ad.IdDiagnostico = d.IdDiagnostico
			INNER JOIN Pacientes p ON a.IdPaciente = p.IdPaciente
			WHERE ad.IdDiagnostico = @IdDiagnostico
			AND a.FechaEgreso IS NOT NULL
			AND a.FyHFinal IS NOT NULL
			AND (a.FechaIngreso BETWEEN DATEADD(DAY, -30, @FechaInicio) AND @FechaFin)
		),
		CTE_Distritos AS (
			-- CTE para contar los distritos únicos afectados en el mes actual y el mes anterior
			SELECT
				di.IdDistrito, 
				di.Nombre AS NombreDistrito, 
				pro.Nombre AS NombreProvincia,
				CASE 
					WHEN af.FechaIngreso >= @FechaInicio AND af.FechaIngreso <= @FechaFin THEN 1
					ELSE 0
				END AS MesActual,
				CASE 
					WHEN af.FechaIngreso >= DATEADD(DAY, -30, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN 1
					ELSE 0
				END AS MesAnterior
			FROM CTE_AtencionesFiltradas af
			INNER JOIN Distritos di ON af.IdDistritoProcedencia = di.IdDistrito
			INNER JOIN Provincias pro ON di.IdProvincia = pro.IdProvincia
		),
		CTE_Resumen AS (
			-- CTE para obtener el total de pacientes y atenciones en ambos meses
			SELECT
				-- Total de pacientes únicos en el mes actual
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= @FechaInicio AND af.FechaIngreso <= @FechaFin THEN af.IdPaciente 
					END) AS TotalPacientesUnicosActual,
				
				-- Total de pacientes únicos en el mes anterior
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= DATEADD(DAY, -30, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN af.IdPaciente 
					END) AS TotalPacientesUnicosMesAnterior,

				-- Total de atenciones en el mes actual
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= @FechaInicio AND af.FechaIngreso <= @FechaFin THEN af.IdAtencion
					END) AS TotalAtencionesMesActual,
				
				-- Total de atenciones en el mes anterior
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= DATEADD(DAY, -30, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN af.IdAtencion
					END) AS TotalAtencionesMesAnterior
			FROM CTE_AtencionesFiltradas af
		)
		SELECT
			-- Datos de distritos
			(SELECT COUNT(DISTINCT IdDistrito) FROM CTE_Distritos WHERE MesActual = 1) AS DistritosAfectadosMesActual,
			(SELECT COUNT(DISTINCT IdDistrito) FROM CTE_Distritos WHERE MesAnterior = 1) AS DistritosAfectadosMesAnterior,
			((SELECT COUNT(DISTINCT IdDistrito) FROM CTE_Distritos WHERE MesActual = 1) - 
			(SELECT COUNT(DISTINCT IdDistrito) FROM CTE_Distritos WHERE MesAnterior = 1)) AS DiferenciaDistritosAfectados,

			-- Datos de pacientes y atenciones
			TotalPacientesUnicosActual,
			TotalPacientesUnicosMesAnterior,
			-- Porcentaje de cambio en pacientes
			CASE 
				WHEN TotalPacientesUnicosMesAnterior > 0 THEN 
					((TotalPacientesUnicosActual - TotalPacientesUnicosMesAnterior) * 100.0) / TotalPacientesUnicosMesAnterior
				ELSE 100.0 
			END AS PorcentajeCambioPacientes,

			TotalAtencionesMesActual,
			TotalAtencionesMesAnterior,
			-- Porcentaje de cambio en atenciones
			CASE 
				WHEN TotalAtencionesMesAnterior > 0 THEN 
					((TotalAtencionesMesActual - TotalAtencionesMesAnterior) * 100.0) / TotalAtencionesMesAnterior
				ELSE 100.0 
			END AS PorcentajeCambioAtenciones,

			-- Ratio de retorno (Atenciones / Pacientes)
			CASE 
				WHEN TotalPacientesUnicosActual > 0 THEN 
					TotalAtencionesMesActual * 1.0 / TotalPacientesUnicosActual
				ELSE 0.0 
			END AS RatioDeRetorno
		FROM CTE_Resumen;
	`

	QUERY_OBTENER_SEXO_POR_DIAGNOSTICO = `
		SELECT 
			CASE 
				WHEN p.IdTipoSexo = 1 THEN 'Masculino'
				WHEN p.IdTipoSexo = 2 THEN 'Femenino'
				ELSE 'Desconocido'
			END AS Sexo,
			COUNT(distinct a.IdPaciente) AS CantidadAtenciones
		FROM Atenciones a
		INNER JOIN Pacientes p ON a.IdPaciente = p.IdPaciente
		INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
		WHERE ad.IdDiagnostico = @IdDiagnostico
		AND a.FechaIngreso >= @FechaInicio
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaEgreso IS NOT NULL
		AND a.FyHFinal IS NOT NULL
		GROUP BY p.IdTipoSexo
		ORDER BY Sexo;
	`

	QUERY_OBTENER_EDADES_POR_DIAGNOSTICO = `
		SELECT 
			CASE 
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'A' AND a.Edad >= 60) THEN '60+'

				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 >= 60) THEN '60+'

				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 >= 60) THEN '60+'

				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 >= 60) THEN '60+'

				ELSE 'Desconocido'
			END AS RangoEdad,
			COUNT(DISTINCT a.IdAtencion) AS CantidadAtenciones
		FROM Atenciones a
		INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
		INNER JOIN TiposEdad te ON a.IdTipoEdad = te.IdTipoEdad
		WHERE a.FechaIngreso >= @FechaInicio
		AND a.FechaIngreso <= @FechaFin
		AND ad.IdDiagnostico = @IdDiagnostico
		GROUP BY 
			CASE 
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'A' AND a.Edad >= 60) THEN '60+'

				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'M' AND a.Edad / 12 >= 60) THEN '60+'

				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'D' AND a.Edad / 365 >= 60) THEN '60+'

				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 0 AND 5) THEN '0-5'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 6 AND 11) THEN '6-11'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 12 AND 17) THEN '12-17'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 18 AND 29) THEN '18-29'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 30 AND 59) THEN '30-59'
				WHEN (te.Codigo = 'H' AND a.Edad / 8760 >= 60) THEN '60+'

				ELSE 'Desconocido'
			END
		ORDER BY CantidadAtenciones DESC;
	`

	QUERY_OBTENER_CANTIDADES_POR_DISTRITO = `
		WITH CTE_AtencionesFiltradas AS (
			SELECT a.IdAtencion, p.IdPaciente, p.IdDistritoProcedencia
			FROM Atenciones a
			INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
			INNER JOIN Diagnosticos d ON ad.IdDiagnostico = d.IdDiagnostico
			INNER JOIN Pacientes p ON a.IdPaciente = p.IdPaciente
			WHERE a.FechaIngreso >= @FechaInicio
			AND a.FechaIngreso <= @FechaFin
			AND a.FechaEgreso IS NOT NULL
			AND a.FyHFinal IS NOT NULL
			AND ad.IdDiagnostico = @IdDiagnostico
		)

		SELECT 
			di.IdDistrito, 
			di.Nombre AS NombreDistrito, 
			pro.Nombre AS NombreProvincia,
			COUNT(DISTINCT af.IdAtencion) AS CantidadAtenciones
		FROM CTE_AtencionesFiltradas af
		INNER JOIN Distritos di ON af.IdDistritoProcedencia = di.IdDistrito
		INNER JOIN Provincias pro ON di.IdProvincia = pro.IdProvincia
		GROUP BY di.IdDistrito, di.Nombre, pro.Nombre
		ORDER BY CantidadAtenciones DESC, di.Nombre, pro.Nombre;
	`

	QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_DIA = `
		SELECT 
			CAST(a.FechaIngreso AS DATE) AS Fecha,
			COUNT(a.IdAtencion) AS CantidadAtenciones
		FROM 
			Atenciones a
		INNER JOIN 
			AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
		INNER JOIN 
			Diagnosticos d ON ad.IdDiagnostico = d.IdDiagnostico
		WHERE 
			a.FechaIngreso >= @FechaInicio
			AND a.FechaIngreso <= @FechaFin
			AND a.FechaEgreso IS NOT NULL
			AND a.FyHFinal IS NOT NULL
			AND ad.IdDiagnostico = @IdDiagnostico
		GROUP BY 
			CAST(a.FechaIngreso AS DATE)
		ORDER BY 
			Fecha;
	`
)
