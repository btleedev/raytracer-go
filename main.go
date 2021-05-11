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

	cameraLookFrom, cameraLookAt, cameraUp, cameraFocusPoint, shapes, lights := sample()
	imageSpec := raytracer.ImageSpec{
		Width:                           imageWidth,
		Height:                          imageHeight,
		AntiAliasingFactor:              antiAliasingFactor,
		RayTracingMaxDepth:              raytracingMaxDepth,
		SoftShadowMonteCarloRepetitions: softShadowMonteCarloRepetitions,

		ImageLocation: "out.png",
	}
	scene := raytracer.Scene{
		CameraLookFrom:   cameraLookFrom,
		CameraLookAt:     cameraLookAt,
		CameraUp:         cameraUp,
		CameraFocusPoint: cameraFocusPoint,
		CameraAperature:  cameraAperature,
		CameraFov:        cameraFovDegrees,
		Shapes:           shapes,
		Lights:           lights,
	}
	raytracer.GenerateImage(imageSpec, scene)
}
