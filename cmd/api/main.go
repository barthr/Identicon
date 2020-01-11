package main

import (
	"log"
	"net/http"
	"os"

	"github.com/barthr/identicon"
)

func getEnvOr(key string, orValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return orValue
	}
	return val
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/identicon/generate", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "no name given", http.StatusPreconditionFailed)
			return
		}

		log.Printf("generating identicon for %s", name)

		w.Header().Set("Content-Type", "image/png")
		if err := identicon.Generate([]byte(name)).WriteImage(w); err != nil {
			http.Error(w, "failed generating identicon", http.StatusInternalServerError)
		}
	})
	if err := http.ListenAndServe(":"+getEnvOr("PORT", "8080"), mux); err != nil {
		log.Fatalf("failed listening to web server because %v", err)
	}
}
