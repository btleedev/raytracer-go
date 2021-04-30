package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"image"
	"image/png"
	"math/rand"
	"os"
)

const antiAliasingFactor = 32
const boundingBoxMaxSize = 1000
const raytracingMaxDepth = 16
const cameraAperature = 0.015
const cameraFovDegrees = 60
const imageWidth = 640  // 3840
const imageHeight = 360 // 2160
const softShadowMonteCarloRepetitions = 32
const softShadowMonteCarloMaxLengthDeviation = 0.25

type imageSpec struct {
	width              int
	height             int
	antiAliasingFactor int
}

type raytraceJob struct {
	i int
	j int
}

type raytraceResult struct {
	pixelIdx       int
	pixelColorFrac r3.Vec
}

func GenerateImage() {
	imageSpec := imageSpec{
		imageWidth,
		imageHeight,
		antiAliasingFactor,
	}
	theScene := tesla(imageSpec)
	myImage := image.NewRGBA(image.Rect(0, 0, imageSpec.width, imageSpec.height))
	jobs := make(chan raytraceJob, imageSpec.height*imageSpec.width)
	results := make(chan raytraceResult, imageSpec.height*imageSpec.width)
	workers := 16
	for i := 0; i < workers; i++ {
		go computePixel(i, &imageSpec, theScene.camera, theScene.shapes, theScene.lights, jobs, results)
	}

	for j := imageSpec.height - 1; j >= 0; j-- {
		for i := 0; i < imageSpec.width; i++ {
			jobs <- raytraceJob{
				i: i,
				j: j,
			}
		}
	}
	close(jobs)

	count := 0
	for j := imageSpec.height - 1; j >= 0; j-- {
		for i := 0; i < imageSpec.width; i++ {
			result := <-results
			myImage.Pix[result.pixelIdx+0] = uint8(result.pixelColorFrac.X * 255.99) // 1st pixel red
			myImage.Pix[result.pixelIdx+1] = uint8(result.pixelColorFrac.Y * 255.99) // 1st pixel green
			myImage.Pix[result.pixelIdx+2] = uint8(result.pixelColorFrac.Z * 255.99) // 1st pixel blue
			myImage.Pix[result.pixelIdx+3] = 255                                     // 1st pixel alpha

			count++
			if count%1000 == 0 {
				fmt.Printf("%.2f%% pixels rendered\n", float64(count)/float64(imageSpec.height*imageSpec.width)*100.0)
			}
		}
	}

	outputFile, err := os.Create("out.png")
	if err != nil {
		panic("failed to create image")
	}
	defer outputFile.Close()
	png.Encode(outputFile, myImage)
}

func computePixel(id int, scene *imageSpec, camera *camera, shapes *[]shape, lights *[]light, jobs <-chan raytraceJob, results chan<- raytraceResult) {
	for job := range jobs {
		pixelColor := r3.Vec{}
		for s := 0; s < scene.antiAliasingFactor; s++ {
			u := (float64(job.i) + rand.Float64()) / float64(scene.width)
			v := (float64(job.j) + rand.Float64()) / float64(scene.height)
			ray := camera.getRay(u, v)
			pixelColor = r3.Add(pixelColor, color(&ray, shapes, lights, 0))
		}
		pixelColor = r3.Scale(1.0/float64(scene.antiAliasingFactor), pixelColor)
		pixelColor = r3.Vec{
			X: pixelColor.X,
			Y: pixelColor.Y,
			Z: pixelColor.Z,
		}
		pixelIdx := (((scene.height - 1 - job.j) * scene.width) + job.i) * 4

		// fmt.Printf("Worker %v finished job for (%v, %v)\n", id, job.i, job.j)
		results <- raytraceResult{
			pixelIdx:       pixelIdx,
			pixelColorFrac: pixelColor,
		}
	}
}

func color(r *ray, shapes *[]shape, lights *[]light, depth int) r3.Vec {
	var hit, minHitRecord = trace(r, shapes, 0.0)
	if hit {
		if depth < raytracingMaxDepth {
			shouldTrace, attenuation, scattered, terminalColor := (*minHitRecord.material).scatter(r, minHitRecord, shapes, lights)
			if shouldTrace {
				recColor := color(&scattered, shapes, lights, depth+1)
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

	// background color
	return r3.Vec{
		X: 0 / 255.0,
		Y: 0 / 255.0,
		Z: 0 / 255.0,
	}
}
