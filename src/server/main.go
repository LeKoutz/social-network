package server

import (
	"fmt"
	"forum/src/models"
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
	var dbPath string = "./db.db"
	var ip, port string
	for i := 1; i < len(args); i++ {
		if args[i] == "--db-path" && i+1 < len(args) {
			dbPath = args[i+1]
			i++
		}
	}
	if err := models.InitDB(dbPath); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	var positionalArgs []string
	for i := 1; i < len(args); i++ {
		if args[i] == "--db-path" {
			i++ // skip the flag value too
		} else {
			positionalArgs = append(positionalArgs, args[i])
		}
	}
	if len(positionalArgs) == 2 {
		ip = positionalArgs[0]
		port = positionalArgs[1]
	} else if len(positionalArgs) == 1 {
		ip = "127.0.0.1"
		port = positionalArgs[0]
	} else if len(positionalArgs) == 0 {
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
