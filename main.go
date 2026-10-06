package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()

	apiCfg := &apiConfig{}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", statusHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.countHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHandler)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func (cfg *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {
	message_string := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(message_string))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
