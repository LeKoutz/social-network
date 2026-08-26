package server

import (
	"forum/src/router"
	"net/http"
)

func startServer(ip, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router.RoutesHandler)
	return http.ListenAndServe(ip+":"+port, mux)
}
