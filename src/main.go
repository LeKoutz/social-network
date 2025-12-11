package forum

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func usage(programName string) {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "./%s [ip] [port]\n", programName)
	fmt.Fprintf(os.Stderr, "./%s [port]\n", programName)
	fmt.Fprintf(os.Stderr, "./%s\n", programName)
}

// Entry point for the program
func Main(args []string) {
	// TODO pass args and populate them maybe, otherwise default
	// maybe initialize something first if needed...
	// for example the database?!?!
	// maybe that could be a flag...
	if err := Init(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// e.g. --init
	var ip, port string
	// in case the flag is not passed, the auto intialization could be triggered
	if len(args) == 3 {
		ip = args[1]
		port = args[2]
	} else if len(args) == 2 {
		ip = "127.0.0.1"
		port = args[1]
	} else if len(args) == 1 {
		ip = "127.0.0.1"
		port = "8080"
	} else {
		programName := filepath.Base(args[0])
		usage(programName)
		os.Exit(1)
	}
	log.Printf("http://%s:%s/", ip, port)
	startServer(ip, port)
}
