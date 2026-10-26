package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/relentlessworks/fuzzkit/internal/api"
	"github.com/relentlessworks/fuzzkit/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.NewHandler(cfg.NoAuth)
	mux := handler.Routes()

	// Wrap with logging middleware
	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: logRequest(mux),
	}

	fmt.Fprintf(os.Stderr, "fuzzkit listening on %s\n", cfg.Addr)
	if cfg.NoAuth {
		fmt.Fprintf(os.Stderr, "WARNING: authentication disabled (dev mode)\n")
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.URL.RawQuery)
		next.ServeHTTP(w, r)
	})
}
