package serve

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zodimo/go-compose/cmd/go-compose/internal/build"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	var (
		httpAddr = fs.String("http", ":8080", "HTTP bind address to serve")
		pkgPath  = "."
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() > 0 {
		pkgPath = fs.Arg(0)
	}

	// Create temporary directory for build output
	tmpDir, err := os.MkdirTemp("", "go-compose-serve-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("Building %s for web...\n", pkgPath)
	// Initial build
	// TODO: Implement on-demand build handler for better DX
	if err := build.BuildJS(tmpDir, pkgPath); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("Serving on %s...\n", *httpAddr)

	// Simple file server for now
	fsHandler := http.FileServer(http.Dir(tmpDir))

	http.Handle("/", fsHandler)

	return http.ListenAndServe(*httpAddr, nil)
}
