// Command api-check enforces the go-compose API-purity contracts.
//
// It inspects the exported signatures of the public package trees
// (compose/, modifiers/, theme/, runtime/, pkg/) and fails when any exported
// symbol references a type from a gioui.org package.
//
// Usage:
//
//	go run ./cmd/api-check          # signature purity of the public trees
//	go run ./cmd/api-check -all     # + engine-import restriction (Phase 3)
//
// Exit codes:
//
//	0 — public API is gioui-free
//	1 — one or more violations found
//	2 — analyzer itself failed (load error, broken tree)
//
// Invoked via `make check-api`.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zodimo/go-compose/internal/api-check"
)

func main() {
	all := flag.Bool("all", false,
		"also enforce the engine-import restriction (Rule B) and extend the signature rule to non-public non-seam packages")
	flag.Parse()

	violations, err := apicheck.Check(apicheck.Options{All: *all})
	if err != nil {
		fmt.Fprintln(os.Stderr, "api-check:", err)
		os.Exit(2)
	}

	for _, v := range violations {
		fmt.Println(v.String())
	}
	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "api-check: %d violation(s) found; public API is not gioui-free\n", len(violations))
		os.Exit(1)
	}
	fmt.Println("api-check: ok — public API is gioui-free")
}
