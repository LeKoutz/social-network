package server

import (
	"forum/src/controllers"
	"net/http"
)

func startServer(ip, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.RoutesHandler)
	// controllers.RegisterRoutes(mux)
	// controllers.Middleware(mux)
	return http.ListenAndServe(ip+":"+port, mux)
}
