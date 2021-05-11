package main

import (
	"example.com/hello/raytracer"
)

const antiAliasingFactor = 32
const cameraAperature = 0.015
const cameraFovDegrees = 60
const imageWidth = 640  // 3840
const imageHeight = 360 // 2160
const raytracingMaxDepth = 16
const softShadowMonteCarloRepetitions = 16

func main() {
	// CPU profiling by default
	// defer profile.Start().Stop()

	imageSpec := raytracer.ImageSpec{
		Width:                           imageWidth,
		Height:                          imageHeight,
		AntiAliasingFactor:              antiAliasingFactor,
		CameraAperature:                 cameraAperature,
		CameraFov:                       cameraFovDegrees,
		RayTracingMaxDepth:              raytracingMaxDepth,
		SoftShadowMonteCarloRepetitions: softShadowMonteCarloRepetitions,
	}

	raytracer.GenerateImage(imageSpec)
}
