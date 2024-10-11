package main

import (
	"log"
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
	"github.com/gskll/chirpy2/internal/handlers"
	"github.com/gskll/chirpy2/internal/middleware"
)

func main() {
	var (
		mux        = http.NewServeMux()
		cfg        = config.NewApiConfig()
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
