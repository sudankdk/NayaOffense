package transport

import (
	"log"
	"net/http"
)

func StartServer(addr string, handler http.Handler) error {
	log.Print("Starting server on port: " + addr)
	return http.ListenAndServe(":"+addr, handler)
}
