package forum

import (
	"net/http"
)

func startServer(ip, port string) error {
	http.HandleFunc("/", routesHandler)
	return http.ListenAndServe(ip+":"+port, nil)
}
