package main

import (
	"github.com/TomasCruz/event-processing/internal/generator"
)

// cmd directory moved out of internal, it is now only invoking commands within internal
func main() {
	generator.Run()
}
