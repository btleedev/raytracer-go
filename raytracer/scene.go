package raytracer

import (
	"fmt"
	"github.com/hschendel/stl"
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
)

type scene struct {
	camera *camera
	shapes *[]shape
	lights *[]light
}

func bunny(is imageSpec) scene {
	return genericShowCase(is, fromStlFile(
		"Bunny-LowPoly.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: -45, Z: 0})
			(*sh).translate(r3.Vec{X: 0, Y: -1, Z: 0})
		},
	))
}

func koala(is imageSpec) scene {
	return genericShowCase(is, fromStlFile(
		"Koala.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: -45, Z: 0})
			(*sh).translate(r3.Vec{X: 0, Y: -1, Z: 0})
		},
	))
}

func genericShowCase(is imageSpec, centerShapes []shape) scene {
	lookFrom := r3.Vec{X: 0, Y: 0, Z: -3}
	lookAt := r3.Vec{X: 0, Y: 0, Z: 0}
	lookFromMinusLookAt := r3.Sub(lookFrom, lookAt)
	cam := NewCamera(
		lookFrom,
		lookAt,
		r3.Vec{X: 0, Y: 1, Z: 0},
		cameraFovDegrees,
		float64(is.width)/float64(is.height),
		cameraAperature,
		math.Sqrt(lookFromMinusLookAt.X*lookFromMinusLookAt.X+lookFromMinusLookAt.Y*lookFromMinusLookAt.Y+lookFromMinusLookAt.Z*lookFromMinusLookAt.Z),
	)
	shapes := append(append(centerShapes, mirrorWalls(3, 0.75)...), floor(-1)...)
	lights := []light{
		ambientLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 241 / 255.0,
				Z: 224 / 255.0,
			},
			lightIntensity: 0.1,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 247 / 255.0,
				Z: 41 / 255.0,
			},
			lightIntensity: 3.0,
			position: r3.Vec{
				X: 0,
				Y: 2.9,
				Z: 0,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: 0, Y: 2.9, Z: 0}),
			angle:     25,
		},
	}
	return scene{
		camera: &cam,
		shapes: &shapes,
		lights: &lights,
	}
}

func floor(yCoord float64) []shape {
	return []shape{
		&triangle{
			pointA:      r3.Vec{X: 10000, Y: yCoord, Z: 10000},
			pointB:      r3.Vec{X: 10000, Y: yCoord, Z: -10000},
			pointC:      r3.Vec{X: -10000, Y: yCoord, Z: 10000},
			singleSided: true,
			mat: phongBlinn{
				specValue:     0.0,
				specShininess: 0.0,
				color:         r3.Vec{X: 255.0 / 255.0, Y: 235.0 / 255.0, Z: 205.0 / 255.0},
			},
		},
		&triangle{
			pointA:      r3.Vec{X: -10000, Y: yCoord, Z: -10000},
			pointB:      r3.Vec{X: -10000, Y: yCoord, Z: 10000},
			pointC:      r3.Vec{X: 10000, Y: yCoord, Z: -10000},
			singleSided: true,
			mat: phongBlinn{
				specValue:     0.0,
				specShininess: 0.0,
				color:         r3.Vec{X: 255.0 / 255.0, Y: 235.0 / 255.0, Z: 205.0 / 255.0},
			},
		},
	}
}

func mirrorWalls(radius float64, albedo float64) []shape {
	return []shape{
		// back
		&triangle{
			pointA:      r3.Vec{X: radius, Y: radius, Z: -radius},
			pointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
			pointC:      r3.Vec{X: radius, Y: -radius, Z: -radius},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		&triangle{
			pointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
			pointB:      r3.Vec{X: radius, Y: -radius, Z: -radius},
			pointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		// right
		&triangle{
			pointA:      r3.Vec{X: radius, Y: radius, Z: 3},
			pointB:      r3.Vec{X: radius, Y: radius, Z: -radius},
			pointC:      r3.Vec{X: radius, Y: -radius, Z: 3},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		&triangle{
			pointA:      r3.Vec{X: radius, Y: -radius, Z: -radius},
			pointB:      r3.Vec{X: radius, Y: -radius, Z: 3},
			pointC:      r3.Vec{X: radius, Y: radius, Z: -radius},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		// left
		&triangle{
			pointA:      r3.Vec{X: -radius, Y: radius, Z: 3},
			pointB:      r3.Vec{X: -radius, Y: -radius, Z: 3},
			pointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		&triangle{
			pointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
			pointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
			pointC:      r3.Vec{X: -radius, Y: -radius, Z: 3},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		// back
		&triangle{
			pointA:      r3.Vec{X: radius, Y: radius, Z: 3},
			pointB:      r3.Vec{X: radius, Y: -radius, Z: 3},
			pointC:      r3.Vec{X: -radius, Y: radius, Z: 3},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
		&triangle{
			pointA:      r3.Vec{X: -radius, Y: -radius, Z: 3},
			pointB:      r3.Vec{X: -radius, Y: radius, Z: 3},
			pointC:      r3.Vec{X: radius, Y: -radius, Z: 3},
			singleSided: true,
			mat:         metal{albedo: r3.Vec{X: albedo, Y: albedo, Z: albedo}},
		},
	}
}

func fromStlFile(stlFileName string, mutator func(shape *shape)) []shape {
	stlFile, err := stl.ReadFile(stlFileName)
	if err != nil {
		panic("failed to load .stl file")
	}

	shapes := make([]shape, len(stlFile.Triangles))
	for i, stlTriangle := range stlFile.Triangles {
		shapes[i] = &triangle{
			pointA:      r3.Vec{X: float64(stlTriangle.Vertices[0][0]), Y: float64(stlTriangle.Vertices[0][1]), Z: float64(stlTriangle.Vertices[0][2])},
			pointB:      r3.Vec{X: float64(stlTriangle.Vertices[1][0]), Y: float64(stlTriangle.Vertices[1][1]), Z: float64(stlTriangle.Vertices[1][2])},
			pointC:      r3.Vec{X: float64(stlTriangle.Vertices[2][0]), Y: float64(stlTriangle.Vertices[2][1]), Z: float64(stlTriangle.Vertices[2][2])},
			singleSided: false,
			mat: phongBlinn{
				specValue:     2,
				specShininess: 5,
				color:         r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()},
			},
		}
		mutator(&shapes[i])
	}

	fmt.Printf("Loaded stl file %s, %v triangles\n", stlFileName, len(shapes))
	return shapes
}
