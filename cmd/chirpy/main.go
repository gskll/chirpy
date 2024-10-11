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
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	var (
		dbQueries  = database.New(db)
		mux        = http.NewServeMux()
		cfg        = config.NewApiConfig(dbQueries)
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

	port := "8080"
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
