package forum

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func usage(programName string) {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "./%s [--db-path <path>] [ip] [port]\n", programName)
	fmt.Fprintf(os.Stderr, "./%s [--db-path <path>] [port]\n", programName)
	fmt.Fprintf(os.Stderr, "./%s [--db-path <path>]\n", programName)
}

// Entry point for the program
func Main(args []string) {
	// TODO pass args and populate them maybe, otherwise default
	// maybe initialize something first if needed...
	// for example the database?!?!
	// maybe that could be a flag...
	var dbPath string = "./data/db.db"
	var ip, port string

	for i := 1; i < len(args); i++ {
		if args[i] == "--db-path" && i+1 < len(args) {
			dbPath = args[i+1]
			i++
		}
	}

	if err := InitDB(dbPath); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// e.g. --init
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
	if err := startServer(ip, port); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
