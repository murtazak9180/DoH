package server

import (
	"doh-server/config"
	"net/http"
)

func handleDNSQuery(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

	}

	if r.Method == "GET" {

	}

	//If not either POST or GET, return 405.
	http.Error(w, "The method must either be POST or GET", http.StatusMethodNotAllowed)

}

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/dns-query", handleDNSQuery)
	return mux
}

func Start(cfg config.Config, handler http.Handler) error {
	return http.ListenAndServeTLS(cfg.Port, cfg.CertPath, cfg.KeyPath, handler)

}
