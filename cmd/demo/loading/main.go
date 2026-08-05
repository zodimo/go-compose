package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("Loading Indicator Demo"),
		runtime.Size(unit.NewDpSize(400, 700)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
