package main

import (
	"fmt"
	forum "forum/src"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <password>", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	fmt.Println(forum.HashPassword(os.Args[1]))
}
