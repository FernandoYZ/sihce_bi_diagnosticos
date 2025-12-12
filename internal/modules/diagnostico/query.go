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
		DECLARE @DiferenciaDias INT = DATEDIFF(DAY, @FechaInicio, @FechaFin);

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
			AND (a.FechaIngreso BETWEEN DATEADD(DAY, -@DiferenciaDias, @FechaInicio) AND @FechaFin)
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
					WHEN af.FechaIngreso >= DATEADD(DAY, -@DiferenciaDias, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN 1
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
					WHEN af.FechaIngreso >= DATEADD(DAY, -@DiferenciaDias, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN af.IdPaciente 
				END) AS TotalPacientesUnicosMesAnterior,

				-- Total de atenciones en el mes actual
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= @FechaInicio AND af.FechaIngreso <= @FechaFin THEN af.IdAtencion
				END) AS TotalAtencionesMesActual,
				
				-- Total de atenciones en el mes anterior
				COUNT(DISTINCT CASE 
					WHEN af.FechaIngreso >= DATEADD(DAY, -@DiferenciaDias, @FechaInicio) AND af.FechaIngreso < @FechaInicio THEN af.IdAtencion
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

	// GRÁFICO DE ATENCIONES POR DÍA
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

	// GRÁFICO DE ATENCIONES POR MES
	QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_MES = `
		SELECT 
			CAST(DATEADD(MONTH, DATEDIFF(MONTH, 0, a.FechaIngreso), 0) AS DATE) AS Fecha,
			COUNT(DISTINCT a.IdAtencion) AS CantidadAtenciones
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
		GROUP BY CAST(DATEADD(MONTH, DATEDIFF(MONTH, 0, a.FechaIngreso), 0) AS DATE)
		ORDER BY Fecha;
	`

	// GRÁFICO DE ATENCIONES POR SEMANAS
	QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_SEMANA = `
		WITH Semanas AS (
			SELECT 
				a.IdAtencion,
				a.FechaIngreso,
				DATEDIFF(DAY, @FechaInicio, a.FechaIngreso) / 7 AS SemanaNumero,
				CAST(DATEADD(DAY, (DATEDIFF(DAY, @FechaInicio, a.FechaIngreso) / 7) * 7, @FechaInicio) AS DATE) AS SemanaInicio
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
		)
		SELECT 
			SemanaInicio,
			COUNT(distinct IdAtencion) AS CantidadAtenciones
		FROM 
			Semanas
		GROUP BY SemanaInicio
		ORDER BY SemanaInicio;
	`

	// GRÁFICO DE ATENCIONES ANUAL
	QUERY_OBTENER_CANTIDAD_ATENCIONES_POR_ANIO = `
		SELECT 
			CAST(DATEADD(YEAR, DATEDIFF(YEAR, 0, a.FechaIngreso), 0) AS DATE) AS AñoInicio,
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
			CAST(DATEADD(YEAR, DATEDIFF(YEAR, 0, a.FechaIngreso), 0) AS DATE)
		ORDER BY AñoInicio;
	`

	QUERY_OBTENER_RANGO_EDADES_SEXO = `
		WITH RangoEdadCalculado AS (
			SELECT 
				a.IdPaciente,
				a.Edad,
				p.IdTipoSexo,
				te.Codigo,
				CASE 
					WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 0 AND 5) THEN '0-05'
					WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 6 AND 11) THEN '06-11'
					WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 12 AND 17) THEN '12-17'
					WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 18 AND 29) THEN '18-29'
					WHEN (te.Codigo = 'A' AND a.Edad BETWEEN 30 AND 59) THEN '30-59'
					WHEN (te.Codigo = 'A' AND a.Edad >= 60) THEN '60+'
					
					WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 0 AND 5) THEN '0-05'
					WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 6 AND 11) THEN '06-11'
					WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 12 AND 17) THEN '12-17'
					WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 18 AND 29) THEN '18-29'
					WHEN (te.Codigo = 'M' AND a.Edad / 12 BETWEEN 30 AND 59) THEN '30-59'
					WHEN (te.Codigo = 'M' AND a.Edad / 12 >= 60) THEN '60+'

					WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 0 AND 5) THEN '0-05'
					WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 6 AND 11) THEN '06-11'
					WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 12 AND 17) THEN '12-17'
					WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 18 AND 29) THEN '18-29'
					WHEN (te.Codigo = 'D' AND a.Edad / 365 BETWEEN 30 AND 59) THEN '30-59'
					WHEN (te.Codigo = 'D' AND a.Edad / 365 >= 60) THEN '60+'

					WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 0 AND 5) THEN '0-05'
					WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 6 AND 11) THEN '06-11'
					WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 12 AND 17) THEN '12-17'
					WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 18 AND 29) THEN '18-29'
					WHEN (te.Codigo = 'H' AND a.Edad / 8760 BETWEEN 30 AND 59) THEN '30-59'
					WHEN (te.Codigo = 'H' AND a.Edad / 8760 >= 60) THEN '60+'
					ELSE 'Desconocido'
				END AS RangoEdad
			FROM Atenciones a
			INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
			INNER JOIN Pacientes p ON a.IdPaciente = p.IdPaciente
			INNER JOIN TiposEdad te ON a.IdTipoEdad = te.IdTipoEdad
			WHERE a.FechaIngreso >= @FechaInicio
			AND a.FechaIngreso <= @FechaFin
			AND ad.IdDiagnostico = @IdDiagnostico
			and a.FechaEgreso IS NOT NULL
			and a.FyHFinal IS NOT NULL

		)

		SELECT 
			RangoEdad,
			COUNT(DISTINCT CASE WHEN IdTipoSexo = 1 THEN IdPaciente END) AS Masculino,
			COUNT(DISTINCT CASE WHEN IdTipoSexo = 2 THEN IdPaciente END) AS Femenino
		FROM RangoEdadCalculado
		GROUP BY RangoEdad
		ORDER BY RangoEdad;
	`

	QUERY_CONDICION_PACIENTE = `
		SELECT 
			CASE 
				WHEN a.IdTipoCondicionAlServicio = 1 THEN 'Nuevo'
				WHEN a.IdTipoCondicionAlServicio = 2 THEN 'Reingreso'
				WHEN a.IdTipoCondicionAlServicio = 3 THEN 'Continuador'
				WHEN a.IdTipoCondicionAlServicio = 4 THEN 'Ausente'
				ELSE 'Desconocido'
			END AS TipoCondicionAlServicio,
			COUNT(*) AS Cantidad
		FROM Atenciones a
		INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
		WHERE a.FyHFinal IS NOT NULL
		AND a.FechaEgreso IS NOT NULL
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaIngreso >= @FechaInicio
		AND ad.IdDiagnostico = @IdDiagnostico
		GROUP BY a.IdTipoCondicionAlServicio
		ORDER BY Cantidad DESC;
	`

	QUERY_CLASIFICACION_DIAGNOSTICO = `
		SELECT 
			CASE 
				WHEN ad.IdSubclasificacionDx = 101 THEN 'Definitivo'
				WHEN ad.IdSubclasificacionDx = 102 THEN 'Presuntivo'
				WHEN ad.IdSubclasificacionDx = 103 THEN 'Repetido'
				ELSE 'Desconocido'
			END AS ClasificacionDiagnostico,
			COUNT(*) AS Cantidad
		FROM Atenciones a
		INNER JOIN AtencionesDiagnosticos ad ON a.IdAtencion = ad.IdAtencion
		WHERE a.FyHFinal IS NOT NULL
		AND a.FechaEgreso IS NOT NULL
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaIngreso >= @FechaInicio
		AND ad.IdDiagnostico = @IdDiagnostico
		GROUP BY ad.IdSubclasificacionDx
		ORDER BY Cantidad DESC;
	`

	QUERY_TABLA_RESUMEN_PACIENTES = `
		-- DECLARE @IdDiagnostico INT = 50795;
		-- DECLARE @FechaInicio DATETIME = '20250401';
		-- DECLARE @FechaFin DATETIME = '20250430';
		-- DECLARE @tamanoPagina INT = 100;
		-- DECLARE @numeroPagina INT = 1;
		-- DECLARE @buscarIdServicio INT = NULL;
		-- DECLARE @buscarIdDistrito int = NULL;
		-- Calcular el offset basado en la página actual
		DECLARE @offset INT = (@numeroPagina - 1) * @tamanoPagina;

		SELECT
			CONVERT(VARCHAR(10), a.FechaIngreso, 23) AS FechaCita,
			ISNULL(p.NroDocumento,'SN') AS NroDocumento,
			p.NroHistoriaClinica,
			(p.PrimerNombre + ' ' + isnull(p.SegundoNombre, '') + ' ' + ISNULL(p.TercerNombre, '') + ' ' + p.ApellidoPaterno + ' ' + p.ApellidoMaterno) as NombrePaciente,
			a.Edad,
			case
				when a.IdTipoEdad = 1 then 'años'
				when a.IdTipoEdad = 2 then 'meses'
				when a.IdTipoEdad = 3 then 'días'
				when a.IdTipoEdad = 4 then 'horas'
			end as TipoEdad,
			ad.IdSubclasificacionDx,
			CASE 
				WHEN ad.IdSubclasificacionDx = 101 THEN 'Definitivo'
				WHEN ad.IdSubclasificacionDx = 102 THEN 'Presuntivo'
				WHEN ad.IdSubclasificacionDx = 103 THEN 'Repetido'
				ELSE 'Desconocido'
			END AS ClasificacionDiagnostico,
			a.IdTipoCondicionAlServicio,
			CASE 
				WHEN a.IdTipoCondicionAlServicio = 1 THEN 'Nuevo'
				WHEN a.IdTipoCondicionAlServicio = 2 THEN 'Reingreso'
				WHEN a.IdTipoCondicionAlServicio = 3 THEN 'Continuador'
				WHEN a.IdTipoCondicionAlServicio = 4 THEN 'Ausente'
				ELSE 'Desconocido'
			END AS CondicionPaciente,
			a.idFuenteFinanciamiento,
			CASE a.idFuenteFinanciamiento
				WHEN 9 THEN 'Estrategia'
				WHEN 5 THEN 'Particular' 
				WHEN 4 THEN 'SOAT'
				WHEN 3 THEN 'SIS'
				WHEN 7 THEN 'SALUDPOL'
				ELSE 'Otro'
			END AS FuenteFinanciamiento,
			s.IdServicio,
			RTRIM(s.Nombre) as Servicio
		from Atenciones a
		inner join AtencionesDiagnosticos ad on a.IdAtencion = ad.IdAtencion
		inner join Pacientes p on a.IdPaciente = p.IdPaciente
		inner join Servicios s on a.IdServicioIngreso = s.IdServicio
		where a.FyHFinal is not null
		and a.FechaEgreso is not null
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaIngreso >= @FechaInicio
		AND ad.IdDiagnostico = @IdDiagnostico
		AND (@buscarIdServicio IS NULL OR a.IdServicioIngreso = @buscarIdServicio)
		AND (@buscarIdDistrito IS NULL OR p.IdDistritoProcedencia = @buscarIdDistrito)
		ORDER BY a.FechaIngreso DESC
		OFFSET @offset ROWS
		FETCH NEXT @tamanoPagina ROWS ONLY
		OPTION (FAST 100, FORCE ORDER);
	`

	QUERY_FILTRAR_SERVICIO = `
		select distinct
			s.IdServicio,
			s.Nombre as Servicio
		from Servicios s
		inner join Atenciones a on s.IdServicio = a.IdServicioIngreso
		inner join AtencionesDiagnosticos ad on a.IdAtencion = ad.IdAtencion
		where a.FyHFinal is not null
		and a.FechaEgreso is not null
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaIngreso >= @FechaInicio
		AND ad.IdDiagnostico = @IdDiagnostico
	`

	QUERY_FILTRAR_DISTRITO = `
		select distinct
			d.IdDistrito,
			d.Nombre as Distrito
		from Distritos d
		inner join Pacientes p on d.IdDistrito = p.IdDistritoProcedencia
		inner join Atenciones a on p.IdPaciente = a.IdPaciente
		inner join AtencionesDiagnosticos ad on a.IdAtencion = ad.IdAtencion
		where a.FyHFinal is not null
		and a.FechaEgreso is not null
		AND a.FechaIngreso <= @FechaFin
		AND a.FechaIngreso >= @FechaInicio
		AND ad.IdDiagnostico = @IdDiagnostico
	` 

	QUERY_ACTIVOS_SEGUIMIENTO = `
	` 

)
