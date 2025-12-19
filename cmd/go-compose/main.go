package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zodimo/go-compose/cmd/go-compose/internal/build"
	"github.com/zodimo/go-compose/cmd/go-compose/internal/serve"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "build":
		if err := build.Run(args); err != nil {
			log.Fatalf("build error: %v", err)
		}
	case "serve":
		if err := serve.Run(args); err != nil {
			log.Fatalf("serve error: %v", err)
		}
	case "help":
		usage()
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: go-compose <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  build    Build the application for a specific target (android, js, desktop)")
	fmt.Println("  serve    Serve the application for web development")
	fmt.Println("  help     Show this help message")
}
