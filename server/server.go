package server

import (
	"doh-server/config"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/miekg/dns"
	mydns "github.com/murtazak9180/DoH/dns"
)

func validatePostRequest(r *http.Request) bool {
	return r.URL.Path == "/dns-query" && r.ContentLength > 0 &&
		strings.ToLower(r.Header.Get("Content-Type")) == "application/dns-message"
}

func validateGetRequest(r *http.Request) bool {
	return r.URL.Path == "/dns-query" && r.URL.Query().Get("dns") != "" &&
		strings.ToLower(r.Header.Get("Content-Type")) == "application/dns-message"
}

func wrapperHandle(body []byte, cfg config.Config, w http.ResponseWriter) []byte {
	msg := new(dns.Msg)
	if err := msg.Unpack(body); err != nil {
		http.Error(w, "Failed to unpack the body", http.StatusBadRequest)
		return nil
	}

	var (
		resp *dns.Msg
		err  error
	)

	if cfg.ResolverMode == "resolver" {
		// TODO
	} else if cfg.ResolverMode == "upstream" {
		resp, err = mydns.UpstreamDNS(msg, cfg.UpstreamDNS)
		if err != nil {
			http.Error(w, "Failed to upstream", http.StatusInternalServerError)
			return nil
		}
	}

	packed, err := resp.Pack()
	if err != nil {
		http.Error(w, "Failed to pack the response", http.StatusInternalServerError)
		return nil
	}
	return packed
}

func handleDNSQuery(w http.ResponseWriter, r *http.Request, cfg config.Config) {
	switch r.Method {
	case "POST":
		if !validatePostRequest(r) {
			http.Error(w, "Invalid POST request", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Failed to read the content", http.StatusBadRequest)
			return
		}
		resp := wrapperHandle(body, cfg, w)
		if resp != nil {
			w.Header().Set("Content-Type", "application/dns-message")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}

	case "GET":
		if !validateGetRequest(r) {
			http.Error(w, "Invalid GET request", http.StatusBadRequest)
			return
		}
		q := r.URL.Query().Get("dns")
		body, err := base64.RawURLEncoding.DecodeString(q)
		if err != nil {
			http.Error(w, "Error decoding base64", http.StatusBadRequest)
			return
		}
		resp := wrapperHandle(body, cfg, w)
		if resp != nil {
			w.Header().Set("Content-Type", "application/dns-message")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}

	default:
		http.Error(w, "The method must either be POST or GET", http.StatusMethodNotAllowed)
	}
}

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/dns-query", func(w http.ResponseWriter, r *http.Request) {
		handleDNSQuery(w, r, cfg)
	})
	return mux
}

func Start(cfg config.Config, handler http.Handler) error {
	return http.ListenAndServeTLS(cfg.Port, cfg.CertPath, cfg.KeyPath, handler)
}
