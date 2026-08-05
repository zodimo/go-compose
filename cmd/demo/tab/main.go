package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("Material 3 Tab Demo"),
		runtime.Size(unit.NewDpSize(600, 400)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
