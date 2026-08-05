package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("Pure Compose"),
		runtime.Size(unit.NewDpSize(1024, 768)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
