package server

import (
	"forum/src/controllers"
	"forum/src/models"
	"net/http"
)

func startServer(ip, port string) error {
	controllers.Hub = models.NewHub()
	go controllers.Hub.Run()
	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.RoutesHandler)
	return http.ListenAndServe(ip+":"+port, mux)
}
