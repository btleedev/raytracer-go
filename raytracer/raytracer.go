package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"image"
	"math"
	"math/rand"
	"time"
)

const bvhCentroidJitterFactor = 0.0000000001
const softShadowMonteCarloMaxLengthDeviation = 0.25

type ImageSpec struct {
	Width                           int
	Height                          int
	AntiAliasingFactor              int
	RayTracingMaxDepth              int
	SoftShadowMonteCarloRepetitions int
	WorkerCount                     int
}

type Scene struct {
	CameraLookFrom   r3.Vec
	CameraLookAt     r3.Vec
	CameraUp         r3.Vec
	CameraFocusPoint r3.Vec // when using camera aperature, point to focus on

	CameraAperature float64
	CameraFov       float64 // in degrees

	Shapes []Shape
	Lights []Light
}

type raytraceJob struct {
	i int
	j int
}

type raytraceResult struct {
	pixelIdx       int
	pixelColorFrac r3.Vec
}

func GenerateImage(imageSpec ImageSpec, scene Scene) *image.RGBA {
	lookFromMinusLookAt := r3.Sub(scene.CameraLookFrom, scene.CameraLookAt)
	cam := NewCamera(
		scene.CameraLookFrom,
		scene.CameraLookAt,
		scene.CameraUp,
		scene.CameraFov,
		float64(imageSpec.Width)/float64(imageSpec.Height),
		scene.CameraAperature,
		math.Sqrt(lookFromMinusLookAt.X*lookFromMinusLookAt.X+lookFromMinusLookAt.Y*lookFromMinusLookAt.Y+lookFromMinusLookAt.Z*lookFromMinusLookAt.Z),
	)
	bvh := NewBoundingVolumeHierarchy(&scene.Shapes)
	myImage := image.NewRGBA(image.Rect(0, 0, imageSpec.Width, imageSpec.Height))
	jobs := make(chan raytraceJob, imageSpec.Height*imageSpec.Width)
	results := make(chan raytraceResult, imageSpec.Height*imageSpec.Width)
	workers := imageSpec.WorkerCount
	for i := 0; i < workers; i++ {
		go computePixel(i, &imageSpec, &cam, bvh, &scene.Lights, jobs, results)
	}

	startTime := time.Now()
	for j := imageSpec.Height - 1; j >= 0; j-- {
		for i := 0; i < imageSpec.Width; i++ {
			jobs <- raytraceJob{
				i: i,
				j: j,
			}
		}
	}
	close(jobs)

	count := 0
	for j := imageSpec.Height - 1; j >= 0; j-- {
		for i := 0; i < imageSpec.Width; i++ {
			result := <-results
			myImage.Pix[result.pixelIdx+0] = uint8(result.pixelColorFrac.X * 255.99) // 1st pixel red
			myImage.Pix[result.pixelIdx+1] = uint8(result.pixelColorFrac.Y * 255.99) // 1st pixel green
			myImage.Pix[result.pixelIdx+2] = uint8(result.pixelColorFrac.Z * 255.99) // 1st pixel blue
			myImage.Pix[result.pixelIdx+3] = 255                                     // 1st pixel alpha

			count++
			if count%1000 == 0 {
				fmt.Printf("%.2f%% pixels rendered, %s\n", float64(count)/float64(imageSpec.Height*imageSpec.Width)*100.0, time.Since(startTime).String())
			}
		}
	}

	fmt.Printf("Finished ray tracing in %s\n", time.Since(startTime).String())
	return myImage
}

func computePixel(id int, is *ImageSpec, camera *camera, bvh *boundingVolumeHierarchy, lights *[]Light, jobs <-chan raytraceJob, results chan<- raytraceResult) {
	for job := range jobs {
		pixelColor := r3.Vec{}
		for s := 0; s < is.AntiAliasingFactor; s++ {
			u := (float64(job.i) + rand.Float64()) / float64(is.Width)
			v := (float64(job.j) + rand.Float64()) / float64(is.Height)
			ray := camera.getRay(u, v)
			pixelColor = r3.Add(pixelColor, color(is, &ray, bvh, lights, 0))
		}
		pixelColor = r3.Scale(1.0/float64(is.AntiAliasingFactor), pixelColor)
		pixelColor = r3.Vec{
			X: pixelColor.X,
			Y: pixelColor.Y,
			Z: pixelColor.Z,
		}
		pixelIdx := (((is.Height - 1 - job.j) * is.Width) + job.i) * 4

		// fmt.Printf("Worker %v finished job for (%v, %v)\n", id, job.i, job.j)
		results <- raytraceResult{
			pixelIdx:       pixelIdx,
			pixelColorFrac: pixelColor,
		}
	}
}

func color(is *ImageSpec, r *ray, bvh *boundingVolumeHierarchy, lights *[]Light, depth int) r3.Vec {
	var hit, minHitRecord = bvh.trace(r, 0.0)
	if hit {
		if depth < is.RayTracingMaxDepth {
			shouldTrace, attenuation, scattered, terminalColor := minHitRecord.material.scatter(is, r, minHitRecord, bvh, lights)
			if shouldTrace {
				recColor := color(is, &scattered, bvh, lights, depth+1)
				return r3.Vec{
					X: attenuation.X * recColor.X,
					Y: attenuation.Y * recColor.Y,
					Z: attenuation.Z * recColor.Z,
				}
			} else {
				return terminalColor
			}
		}
	}

	// background Color
	return r3.Vec{
		X: 0 / 255.0,
		Y: 0 / 255.0,
		Z: 0 / 255.0,
	}
}
