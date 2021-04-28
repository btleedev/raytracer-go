package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type scene struct {
	width int
	height int
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
	scene := scene{
		3840,
		2160,
		//600,
		//400,
		100,
	}
	lookFrom := r3.Vec{ X: 3, Y: 1.5, Z: 1.75 }
	lookAt := r3.Vec{ X: 0, Y: 1, Z: 0 }
	lookFromMinusLookAt := r3.Sub(lookFrom, lookAt)
	camera := NewCamera(
		lookFrom,
		lookAt,
		r3.Vec{ X: 0, Y: 1, Z: 0 },
		60,
		float64(scene.width) / float64(scene.height),
		0.0,
		math.Sqrt(lookFromMinusLookAt.X*lookFromMinusLookAt.X + lookFromMinusLookAt.Y*lookFromMinusLookAt.Y + lookFromMinusLookAt.Z*lookFromMinusLookAt.Z),
	)
	shapes := originalShapes()
	lights := []light{
		//ambientLight{
		//	colorFrac: r3.Vec{
		//		X: 0 / 255.0,
		//		Y: 100 / 255.0,
		//		Z: 125 / 255.0,
		//	},
		//	lightIntensity: 0.5,
		//},
		//pointLight{
		//	colorFrac: r3.Vec{
		//		X: 255 / 255.0,
		//		Y: 255 / 255.0,
		//		Z: 255 / 255.0,
		//	},
		//	lightIntensity: 1.0,
		//	position: r3.Vec{
		//		X: 1,
		//		Y: 2,
		//		Z: -2,
		//	},
		//},
		spotLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 255 / 255.0,
				Z: 255 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 0,
				Y: 3,
				Z: 0,
			},
			direction: r3.Sub(r3.Vec{ X: 0, Y: 0, Z: 0 }, r3.Vec{ X: 0, Y: 3, Z: 0 }),
			angle: 40,
		},
	}
	myImage := image.NewRGBA(image.Rect(0, 0, scene.width, scene.height))
	jobs := make(chan raytraceJob, scene.height * scene.width)
	results := make(chan raytraceResult, scene.height * scene.width)
	workers := 16
	for i := 0; i < workers; i++ {
		go computePixel(i, &scene, &camera, &shapes, &lights, jobs, results)
	}

	for j := scene.height-1; j >= 0; j-- {
		for i := 0; i < scene.width; i++ {
			jobs <- raytraceJob{
				i: i,
				j: j,
			}
		}
	}
	close(jobs)

	count := 0
	for j := scene.height-1; j >= 0; j-- {
		for i := 0; i < scene.width; i++ {
			result := <- results
			myImage.Pix[result.pixelIdx+0] = uint8(result.pixelColorFrac.X * 255.99) // 1st pixel red
			myImage.Pix[result.pixelIdx+1] = uint8(result.pixelColorFrac.Y * 255.99) // 1st pixel green
			myImage.Pix[result.pixelIdx+2] = uint8(result.pixelColorFrac.Z * 255.99) // 1st pixel blue
			myImage.Pix[result.pixelIdx+3] = 255                                     // 1st pixel alpha

			count++
			if count % 10000 == 0 {
				fmt.Printf("%v%% done\n", float64(count) / float64(scene.height * scene.width) * 100.0)
			}
		}
	}

	outputFile, err := os.Create("test.png")
	if err != nil {
		panic("failed to create image")
	}
	defer outputFile.Close()
	png.Encode(outputFile, myImage)
}

func computePixel(id int, scene *scene, camera *camera, shapes *[]shape, lights *[]light, jobs <-chan raytraceJob, results chan<- raytraceResult) {
	for job := range jobs {
		pixelColor := r3.Vec{}
		for s := 0; s < scene.antiAliasingFactor; s++ {
			u := (float64(job.i) + rand.Float64()) / float64(scene.width)
			v := (float64(job.j) + rand.Float64()) / float64(scene.height)
			ray := camera.getRay(u, v)
			pixelColor = r3.Add(pixelColor, color(&ray, shapes, lights, 0))
		}
		pixelColor = r3.Scale(1.0 / float64(scene.antiAliasingFactor), pixelColor)
		pixelColor = r3.Vec{
			X: math.Sqrt(pixelColor.X),
			Y: math.Sqrt(pixelColor.Y),
			Z: math.Sqrt(pixelColor.Z),
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
		if depth < 64 {
			shouldTrace, attenuation, scattered, terminalColor := (*minHitRecord.material).scatter(r, minHitRecord, shapes, lights)
			if shouldTrace {
				recColor := color(&scattered, shapes, lights, depth + 1)
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

func randomShapes() []shape {
	n := 22 * 22 + 1 + 3
	shapes := make([]shape, n)
	shapes[0] = sphere{
		center: r3.Vec{ X: 0, Y: -1000, Z: -1 },
		radius: 1000,
		mat: diffuse{ albedo: r3.Vec{ X: 0.5, Y: 0.5, Z: 0.5 } },
	}
	i := 1
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := r3.Vec {}
			for {
				center = r3.Vec{ X: float64(a) + 0.9 * rand.Float64(), Y: 0.2, Z: float64(b) + 0.9 * rand.Float64() }
				centerMinusCenterBalls := r3.Sub(center, r3.Vec{ X: 4, Y: 0.2, Z: 0 })
				if centerMinusCenterBalls.X*centerMinusCenterBalls.X + centerMinusCenterBalls.Y*centerMinusCenterBalls.Y + centerMinusCenterBalls.Z*centerMinusCenterBalls.Z > 0.9 {
					break
				}
			}
			chooseMat := rand.Float64()
			if chooseMat < 0.8 {
				shapes[i] = sphere{
					center: center,
					radius: 0.2,
					mat: diffuse{ albedo: r3.Vec{
						X: rand.Float64()*rand.Float64(),
						Y: rand.Float64()*rand.Float64(),
						Z: rand.Float64()*rand.Float64(),
					} },
				}
			} else if chooseMat < 0.95 {
				shapes[i] = sphere{
					center: center,
					radius: 0.2,
					mat: metal{
						albedo: r3.Vec{
							X: 0.5 * (1 + rand.Float64()),
							Y: 0.5 * (1 + rand.Float64()),
							Z: 0.5 * (1 + rand.Float64()),
						},
						fuzz: 0.5 * rand.Float64(),
					},
				}
			} else {
				shapes[i] = sphere{
					center: center,
					radius: 0.2,
					mat: dielectric{ refractiveIndex: 1.5 },
				}
			}
			i++
		}
	}

	shapes[i] = sphere{
		center: r3.Vec{
			X: 0,
			Y: 1,
			Z: 0,
		},
		radius: 1,
		mat: dielectric{
			refractiveIndex: 1.5,
		},
	}
	i++

	shapes[i] = sphere{
		center: r3.Vec{
			X: -4,
			Y: 1,
			Z: 0,
		},
		radius: 1,
		mat: diffuse{ albedo: r3.Vec{
			X: 0.4,
			Y: 0.2,
			Z: 0.1,
		} },
	}
	i++

	shapes[i] = sphere{
		center: r3.Vec{
			X: 4,
			Y: 1,
			Z: 0,
		},
		radius: 1,
		mat: metal{
			albedo: r3.Vec{
				X: 0.7,
				Y: 0.6,
				Z: 0.5,
			},
			fuzz: 0.0,
		},
	}
	return shapes
}

func originalShapes() []shape {
	return []shape{

		triangle{
			pointA: r3.Vec{ X: 10000, Y: 0, Z: 10000 },
			pointB: r3.Vec{ X: 10000, Y: 0, Z: -10000 },
			pointC: r3.Vec{ X: -10000, Y: 0, Z: 10000 },
			singleSided: true,
			mat: phongBlinn{
				specValue: 0.0,
				specShininess: 0.0,
				color: r3.Vec{ X: 255.0 / 255.0, Y: 235.0 / 255.0, Z: 205.0 / 255.0 },
			},
		},
		triangle{
			pointA: r3.Vec{ X: -10000, Y: 0, Z: -10000 },
			pointB: r3.Vec{ X: -10000, Y: 0, Z: 10000 },
			pointC: r3.Vec{ X: 10000, Y: 0, Z: -10000 },
			singleSided: true,
			mat: phongBlinn{
				specValue: 0.0,
				specShininess: 0.0,
				color: r3.Vec{ X: 255.0 / 255.0, Y: 235.0 / 255.0, Z: 205.0 / 255.0 },
			},
		},
		sphere{
			center: r3.Vec{ X: 0, Y: 0.5, Z: 0 },
			radius: 0.5,
			// mat: diffuse{ albedo: r3.Vec{ X: 0.8, Y: 0.3, Z: 0.3 } },
			mat: phongBlinn{
				specValue: 5.0,
				specShininess: 32.0,
				color: r3.Vec{ X: 1, Y: 0, Z: 0 },
			},
		},
		sphere{
			center: r3.Vec{ X: 1, Y: 0.5, Z: 0 },
			radius: 0.5,
			// mat: metal{ albedo: r3.Vec{ X: 0.8, Y: 0.6, Z: 0.2 } },
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		sphere{
			center: r3.Vec{ X: -1, Y: 0.5, Z: 0 },
			radius: 0.5,
			mat: dielectric{ refractiveIndex: 2.417 },
		},
		// back
		triangle{
			pointA: r3.Vec{ X: 3, Y: 3, Z: -3 },
			pointB: r3.Vec{ X: -3, Y: 3, Z: -3 },
			pointC: r3.Vec{ X: 3, Y: 0, Z: -3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		triangle{
			pointA: r3.Vec{ X: -3, Y: 0, Z: -3 },
			pointB: r3.Vec{ X: 3, Y: 0, Z: -3 },
			pointC: r3.Vec{ X: -3, Y: 3, Z: -3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		// right
		triangle{
			pointA: r3.Vec{ X: 3, Y: 3, Z: 3 },
			pointB: r3.Vec{ X: 3, Y: 3, Z: -3 },
			pointC: r3.Vec{ X: 3, Y: 0, Z: 3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		triangle{
			pointA: r3.Vec{ X: 3, Y: 0, Z: -3 },
			pointB: r3.Vec{ X: 3, Y: 0, Z: 3 },
			pointC: r3.Vec{ X: 3, Y: 3, Z: -3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		// left
		triangle{
			pointA: r3.Vec{ X: -3, Y: 3, Z: 3 },
			pointB: r3.Vec{ X: -3, Y: 0, Z: 3 },
			pointC: r3.Vec{ X: -3, Y: 3, Z: -3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		triangle{
			pointA: r3.Vec{ X: -3, Y: 0, Z: -3 },
			pointB: r3.Vec{ X: -3, Y: 3, Z: -3 },
			pointC: r3.Vec{ X: -3, Y: 0, Z: 3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		// back
		triangle{
			pointA: r3.Vec{ X: 3, Y: 3, Z: 3 },
			pointB: r3.Vec{ X: 3, Y: 0, Z: 3 },
			pointC: r3.Vec{ X: -3, Y: 3, Z: 3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
		triangle{
			pointA: r3.Vec{ X: -3, Y: 0, Z: 3 },
			pointB: r3.Vec{ X: -3, Y: 3, Z: 3 },
			pointC: r3.Vec{ X: 3, Y: 0, Z: 3 },
			singleSided: true,
			mat: metal{ albedo: r3.Vec{ X: 1, Y: 1, Z: 1 } },
		},
	}
}