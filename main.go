// Command airy-wave computes Airy (linear) wave properties and serves a web
// console plus a JSON endpoint backed by the wave model.
//
// Usage:
//
//	go run . -http :8080            # serve the web console + /api on :8080
//	go run . -http :8080 -example-dir example
//
// The web console posts to /api/wave and plots the resulting velocity/pressure
// profile across the water column. All computation is local.
package main

import (
	"flag"
	"fmt"
	"os"

	"airy-wave/internal/server"
)

func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address for the web console and /api")
	exampleDir := flag.String("example-dir", "example", "directory containing example JSON (served at /example/)")
	flag.Parse()

	srv := server.NewServer(*exampleDir)
	fmt.Fprintf(os.Stderr, "airy-wave: listening on %s (examples in %q)\n", *httpAddr, *exampleDir)
	if err := srv.ListenAndServe(*httpAddr); err != nil {
		fmt.Fprintf(os.Stderr, "airy-wave: server error: %v\n", err)
		os.Exit(1)
	}
}
