package server

import (
	"forum/src/controllers"
	"net/http"
)

func startServer(ip, port string) error {
	http.HandleFunc("/", controllers.RoutesHandler)
	return http.ListenAndServe(ip+":"+port, nil)
}
