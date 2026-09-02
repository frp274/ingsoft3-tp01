package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB inicializa la conexión a MySQL, configura el pool de conexiones,
// ejecuta la creación de tablas y siembra los datos iniciales.
func InitDB(dsn string) (*sql.DB, error) {
	if !strings.Contains(dsn, "parseTime=") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión con la base de datos: %w", err)
	}

	// Configuración del connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Reintentos de conexión (espera a que MySQL esté 100% listo)
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		log.Printf("Esperando disponibilidad de MySQL (intento %d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		return nil, fmt.Errorf("no se pudo conectar con MySQL tras varios intentos: %w", pingErr)
	}

	log.Println("Conexión con MySQL establecida correctamente.")

	// Crear tablas si no existen
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error al crear tablas: %w", err)
	}

	// Ejecutar script de seed
	if err := Seed(db); err != nil {
		return nil, fmt.Errorf("error al ejecutar seed: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS usuarios (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			contrasena VARCHAR(255) NOT NULL,
			nombre VARCHAR(100) NOT NULL,
			apellido VARCHAR(100) NOT NULL,
			ciudad VARCHAR(100) NOT NULL,
			edad INT NOT NULL,
			etiquetaFavorita VARCHAR(100) NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS etiquetas (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL UNIQUE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS planes (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nombre VARCHAR(255) NOT NULL,
			idOrganizador INT NOT NULL,
			fechaInicio DATETIME NOT NULL,
			fechaFin DATETIME NOT NULL,
			visibilidad ENUM('publico', 'privado') NOT NULL DEFAULT 'publico',
			descripcion TEXT NULL,
			idEtiqueta INT NOT NULL,
			ubicacion VARCHAR(255) NOT NULL,
			capacidad INT NOT NULL,
			CONSTRAINT fk_planes_organizador FOREIGN KEY (idOrganizador) REFERENCES usuarios(id) ON DELETE CASCADE,
			CONSTRAINT fk_planes_etiqueta FOREIGN KEY (idEtiqueta) REFERENCES etiquetas(id) ON DELETE RESTRICT
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS invitaciones (
			idInvitacion INT AUTO_INCREMENT PRIMARY KEY,
			idPlan INT NOT NULL,
			idUsuario INT NOT NULL,
			estado ENUM('pendiente', 'aceptada', 'rechazada') NOT NULL DEFAULT 'pendiente',
			CONSTRAINT fk_invitaciones_plan FOREIGN KEY (idPlan) REFERENCES planes(id) ON DELETE CASCADE,
			CONSTRAINT fk_invitaciones_usuario FOREIGN KEY (idUsuario) REFERENCES usuarios(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	log.Println("Esquema de tablas validado/creado con éxito.")
	return nil
}

// Seed siembra la base de datos con datos iniciales si las tablas están vacías.
func Seed(db *sql.DB) error {
	// 1. Verificar si ya existen usuarios
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		log.Println("La base de datos ya contiene datos. Saltando paso de Seed.")
		return nil
	}

	log.Println("Ejecutando Seed de datos iniciales...")

	// 2. Insertar Usuario de prueba
	userQuery := `INSERT INTO usuarios (email, contrasena, nombre, apellido, ciudad, edad, etiquetaFavorita)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	userRes, err := db.Exec(userQuery, "juan.perez@example.com", "password123", "Juan", "Pérez", "Córdoba", 26, "Outdoor")
	if err != nil {
		return fmt.Errorf("error al sembrar usuario: %w", err)
	}
	userId, err := userRes.LastInsertId()
	if err != nil {
		return err
	}

	// 3. Insertar 3 Etiquetas base
	tagQuery := `INSERT INTO etiquetas (nombre) VALUES (?), (?), (?)`
	if _, err := db.Exec(tagQuery, "Outdoor", "Fiesta", "Cine"); err != nil {
		return fmt.Errorf("error al sembrar etiquetas: %w", err)
	}

	// Obtener IDs de las etiquetas
	var idOutdoor, idFiesta, idCine int
	if err := db.QueryRow("SELECT id FROM etiquetas WHERE nombre = 'Outdoor'").Scan(&idOutdoor); err != nil {
		return err
	}
	if err := db.QueryRow("SELECT id FROM etiquetas WHERE nombre = 'Fiesta'").Scan(&idFiesta); err != nil {
		return err
	}
	if err := db.QueryRow("SELECT id FROM etiquetas WHERE nombre = 'Cine'").Scan(&idCine); err != nil {
		return err
	}

	// 4. Insertar 5 Planes públicos asociados al usuario y etiquetas
	now := time.Now()
	planes := []struct {
		nombre      string
		idOrg       int64
		fechaInicio time.Time
		fechaFin    time.Time
		visibilidad string
		descripcion string
		idEtiqueta  int
		ubicacion   string
		capacidad   int
	}{
		{
			nombre:      "Trekking al Cerro Champaquí",
			idOrg:       userId,
			fechaInicio: now.Add(48 * time.Hour),
			fechaFin:    now.Add(56 * time.Hour),
			visibilidad: "publico",
			descripcion: "Excursión de senderismo y aventura en las altas cumbres con guía profesional.",
			idEtiqueta:  idOutdoor,
			ubicacion:   "Villa Alpina, Córdoba",
			capacidad:   15,
		},
		{
			nombre:      "Noche Electrónica & Juegos de Mesa",
			idOrg:       userId,
			fechaInicio: now.Add(72 * time.Hour),
			fechaFin:    now.Add(78 * time.Hour),
			visibilidad: "publico",
			descripcion: "Juntada para compartir juegos de mesa modernos, música en vivo y tragos de autor.",
			idEtiqueta:  idFiesta,
			ubicacion:   "Bar Cultural Güemes, Córdoba",
			capacidad:   30,
		},
		{
			nombre:      "Maratón de Cine de Ciencia Ficción",
			idOrg:       userId,
			fechaInicio: now.Add(96 * time.Hour),
			fechaFin:    now.Add(102 * time.Hour),
			visibilidad: "publico",
			descripcion: "Proyección al aire libre de clásicos de Sci-Fi con debate y pochoclos.",
			idEtiqueta:  idCine,
			ubicacion:   "Parque Sarmiento, Córdoba",
			capacidad:   45,
		},
		{
			nombre:      "Kayak y Picnic en el Dique San Roque",
			idOrg:       userId,
			fechaInicio: now.Add(120 * time.Hour),
			fechaFin:    now.Add(126 * time.Hour),
			visibilidad: "publico",
			descripcion: "Tarde de navegación en kayak por el lago y merienda compartida al atardecer.",
			idEtiqueta:  idOutdoor,
			ubicacion:   "Club Náutico, Villa Carlos Paz",
			capacidad:   20,
		},
		{
			nombre:      "Ciclo de Cine Independiente & Cortometrajes",
			idOrg:       userId,
			fechaInicio: now.Add(144 * time.Hour),
			fechaFin:    now.Add(148 * time.Hour),
			visibilidad: "publico",
			descripcion: "Muestra de directores locales y charla abierta con los creadores audiovisuales.",
			idEtiqueta:  idCine,
			ubicacion:   "Cineclub Municipal, Córdoba",
			capacidad:   50,
		},
	}

	planQuery := `INSERT INTO planes (nombre, idOrganizador, fechaInicio, fechaFin, visibilidad, descripcion, idEtiqueta, ubicacion, capacidad)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, p := range planes {
		if _, err := db.Exec(planQuery, p.nombre, p.idOrg, p.fechaInicio, p.fechaFin, p.visibilidad, p.descripcion, p.idEtiqueta, p.ubicacion, p.capacidad); err != nil {
			return fmt.Errorf("error al sembrar plan '%s': %w", p.nombre, err)
		}
	}

	log.Println("Seed de datos completado exitosamente (1 Usuario, 3 Etiquetas, 5 Planes).")
	return nil
}