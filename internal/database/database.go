package database

import (
	"database/sql"
	"fmt"
	"log"
	"sihce_diagnosticos/internal/config"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func ConectarDB(configuracion config.Config) (*sql.DB, error) {
	cadenaConexion := fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=%s;TrustServerCertificate=%s;connection timeout=%d",
		configuracion.DBServer,
		configuracion.DBPort,
		configuracion.DBName,
		configuracion.DBUser,
		configuracion.DBPassword,
		configuracion.DBEncrypt,
		configuracion.DBTrustServerCertificate,
		120,
	)

	// Abrir conexión a la base de datos
	db, err := sql.Open("sqlserver", cadenaConexion)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexion: %v", err)
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("no se puede conectar a la base de datos: %v", err)
	}

	// Configuración de la conexión
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("✓ Conexión exitosa a la base de datos")
	return db, nil
}

func CerrarConexion(db *sql.DB) error {
	// Intentar cerrar la conexión
	err := db.Close()
	if err != nil {
		log.Printf("❌ Error al cerrar la conexión a la base de datos: %v", err)
		return fmt.Errorf("error al cerrar la conexión: %v", err)
	}

	log.Println("✓ Conexión cerrada correctamente")
	return nil
}
