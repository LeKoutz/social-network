package main

import (
	"fmt"
	"os"
	"path/filepath"

	"forum/src/controllers"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <password>", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	fmt.Println(controllers.HashPassword(os.Args[1]))
}
