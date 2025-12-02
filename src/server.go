package forum

import (
	"log"
	"net/http"
)

func startServer(ip, port string) {
	http.HandleFunc("/", routesHandler)
	err := http.ListenAndServe(ip+":"+port, nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
}
