package forum

import (
	"log"
	"net/http"
)

func startServer(ip, port string) {
	http.HandleFunc("/", routesHandler)
	log.Fatal(http.ListenAndServe(ip+":"+port, nil))
	}
