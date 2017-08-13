package web

import (
	"github.com/hemtjanst/hemtjanst/device"
	"log"
	"net/http"
	"time"
)

// Serve serves the webinterface
func Serve(d *device.Manager) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler(d))
	h := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Fatal(h.ListenAndServe())
}
