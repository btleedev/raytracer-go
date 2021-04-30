package main

import (
	"example.com/hello/raytracer"
)

func main() {
	// CPU profiling by default
	// defer profile.Start().Stop()

	raytracer.GenerateImage()
}
