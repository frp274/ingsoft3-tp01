package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/frp274/ingsoft3-tp01/backend/internal/handlers"
	"github.com/frp274/ingsoft3-tp01/backend/internal/repository"
	"github.com/frp274/ingsoft3-tp01/backend/internal/services"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Database string `json:"database"`
}

func main() {
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Printf("Sentry initialization failed: %v", err)
		} else {
			defer sentry.Flush(2 * time.Second)
			log.Println("Sentry inicializado correctamente en backend")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		dbDsn = "root:secret_password_here@tcp(localhost:3306)/planes_db"
	}

	log.Printf("Iniciando conexión con base de datos: %s", dbDsn)
	db, err := repository.InitDB(dbDsn)
	if err != nil {
		log.Fatalf("Error crítico al inicializar la base de datos: %v", err)
	}
	defer db.Close()

	// Inyección de dependencias (Repository -> Service -> Handler)
	usuarioRepo := repository.NewUsuarioRepository(db)
	etiquetaRepo := repository.NewEtiquetaRepository(db)
	planRepo := repository.NewPlanRepository(db)
	invitacionRepo := repository.NewInvitacionRepository(db)

	planService := services.NewPlanService(planRepo)
	planHandler := handlers.NewPlanHandler(planService)

	invitacionService := services.NewInvitacionService(invitacionRepo, planRepo)
	invitacionHandler := handlers.NewInvitacionHandler(invitacionService)

	mux := http.NewServeMux()

	// Endpoint Healthcheck
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:   "ok",
			Message:  "Backend API is healthy and running",
			Database: "connected",
		})
	})

	// Endpoint Planes (GET /api/planes, POST /api/planes)
	mux.HandleFunc("/api/planes", planHandler.HandlePlanes)

	// Endpoint Invitaciones y Mis Planes (US-5)
	mux.HandleFunc("/api/invitaciones", invitacionHandler.HandleInvitaciones)
	mux.HandleFunc("/api/mis-planes", invitacionHandler.HandleMisPlanes)

	// Endpoint listar etiquetas
	mux.HandleFunc("/api/etiquetas", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		etiquetas, err := etiquetaRepo.List(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(etiquetas)
	})

	// Endpoint listar usuarios
	mux.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		usuarios, err := usuarioRepo.List(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(usuarios)
	})

	// Endpoint de prueba de falla para observabilidad TP9
	mux.HandleFunc("/api/falla", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("Falla de prueba del TP9")
		sentry.CaptureException(err)
		sentry.Flush(2 * time.Second)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bienvenido a la API de Planes",
		})
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
