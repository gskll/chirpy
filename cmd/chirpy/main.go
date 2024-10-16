package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gskll/chirpy2/internal/config"
	"github.com/gskll/chirpy2/internal/database"
	"github.com/gskll/chirpy2/internal/handlers"
	"github.com/gskll/chirpy2/internal/middleware"
)

func main() {
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	platform := os.Getenv("PLATFORM")
	dbUrl := os.Getenv("DB_URL")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	var (
		dbQueries  = database.New(db)
		mux        = http.NewServeMux()
		cfg        = config.NewApiConfig(dbQueries, platform, jwtSecret, polkaKey)
		middleware = middleware.NewMiddleware(cfg)
	)

	var (
		filepathRoot      = "./public"
		fileServer        = http.FileServer(http.Dir(filepathRoot))
		fileServerHandler = http.StripPrefix("/app", fileServer)
	)

	mux.Handle("/app/", middleware.Metrics(fileServerHandler))

	handlers.RegisterAdminHandlers("/admin", cfg, mux)
	handlers.RegisterAPIHandlers("/api", cfg, mux)

	wrappedMux := middleware.Logger(mux)

	port := "8080"
	srv := &http.Server{
		Handler: wrappedMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
