package main

import (
	"example.com/hello/raytracer"
	"image/png"
	"os"
)

func main() {
	// CPU profiling by default
	// defer profile.Start().Stop()

	imageLocation := "out.png"
	imageSpec, scene := raytracer.ExampleRegression(640, 380, "./")
	myImage := raytracer.GenerateImage(imageSpec, scene)

	outputFile, err := os.Create(imageLocation)
	if err != nil {
		panic("failed to create image")
	}
	defer outputFile.Close()
	png.Encode(outputFile, myImage)
}
