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
)
