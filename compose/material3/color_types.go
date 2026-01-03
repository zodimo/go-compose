package material3

import "github.com/zodimo/go-compose/compose/ui/graphics"

// to be used instead of passing composer as parameter to every function
type LocalColorReciever = func(c Composer) graphics.Color
