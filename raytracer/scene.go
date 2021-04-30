package raytracer

import (
	"fmt"
	"github.com/hschendel/stl"
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
	"sort"
)

type scene struct {
	camera *camera
	shapes *[]shape
	lights *[]light
}

func bunny(is imageSpec) scene {
	return genericShowCase(is, true, fromStlFile(
		"Istareyn/low-poly-stanford-bunny/Bunny-LowPoly.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: -45, Z: 0})
			(*sh).translate(r3.Vec{X: 0, Y: -1, Z: 0})
		},
	))
}

func bellsprout(is imageSpec) scene {
	return genericShowCase(is, true, fromStlFile(
		"Philin_theBlank/bellsprout-with-flower-pot/Bellsprout_in_Flower_Pot.stl",
		func(sh *shape) {
			(*sh).scale(0.25)
			(*sh).rotate(r3.Vec{X: -90, Y: 165, Z: 0})
			(*sh).translate(r3.Vec{X: 0.20, Y: -0.5, Z: -1})
		},
	))
}

func koala(is imageSpec) scene {
	return genericShowCase(is, false, fromStlFile(
		"TroySlatton/lyman-from-animal-crossing/Lyman.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: 0, Y: 180, Z: 0})
			(*sh).translate(r3.Vec{X: 0.5, Y: -1, Z: 0})
		},
	))
}

func tesla(is imageSpec) scene {
	return genericShowCase(is, true, fromStlFile(
		"Sim3D_/tesla-model-3-for-3d-printing/solid/Tesla Model 3.STL",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: 240, Z: 0})
			(*sh).translate(r3.Vec{X: 0.20, Y: -1, Z: 0})
		},
	))
}

func genericShowCase(is imageSpec, withMirrors bool, centerShapes []shape) scene {
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
	shapes := append(centerShapes, floor(-1)...)
	if withMirrors {
		shapes = append(shapes, mirrorWalls(3, 0.75)...)
	}
	lights := []light{
		//ambientLight{
		//	colorFrac: r3.Vec{
		//		X: 255 / 255.0,
		//		Y: 241 / 255.0,
		//		Z: 224 / 255.0,
		//	},
		//	lightIntensity: 0.1,
		//},
		spotLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 255 / 255.0,
				Z: 255 / 255.0,
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

	// small optimization, sort based off distance of Z axis from origin so our bounding boxes are a little better
	sort.SliceStable(stlFile.Triangles, func(i, j int) bool {
		z11 := float64(stlFile.Triangles[i].Vertices[0][2])
		z12 := float64(stlFile.Triangles[i].Vertices[1][2])
		z13 := float64(stlFile.Triangles[i].Vertices[2][2])
		z21 := float64(stlFile.Triangles[i].Vertices[0][2])
		z22 := float64(stlFile.Triangles[i].Vertices[1][2])
		z23 := float64(stlFile.Triangles[i].Vertices[2][2])
		return math.Max(z11, math.Max(z12, z13)) < math.Max(z21, math.Max(z22, z23))
	})

	pMinDefault := r3.Vec{
		X: math.MaxInt64,
		Y: math.MaxInt64,
		Z: math.MaxInt64,
	}
	pMaxDefault := r3.Vec{
		X: math.MinInt64,
		Y: math.MinInt64,
		Z: math.MinInt64,
	}
	pMin := pMinDefault
	pMax := pMaxDefault
	shapes := make([]shape, 0, (len(stlFile.Triangles)/boundingBoxMaxSize)+1)
	boundingBoxShapes := make([]shape, 0, boundingBoxMaxSize)
	for i, stlTriangle := range stlFile.Triangles {
		idx := i % boundingBoxMaxSize
		s := triangle{
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
		boundingBoxShapes = append(boundingBoxShapes, &s)
		mutator(&boundingBoxShapes[idx])

		// calculate bounding box
		pMin.X = math.Min(pMin.X, math.Min(s.pointA.X, math.Min(s.pointB.X, s.pointC.X)))
		pMin.Y = math.Min(pMin.Y, math.Min(s.pointA.Y, math.Min(s.pointB.Y, s.pointC.Y)))
		pMin.Z = math.Min(pMin.Z, math.Min(s.pointA.Z, math.Min(s.pointB.Z, s.pointC.Z)))
		pMax.X = math.Max(pMax.X, math.Max(s.pointA.X, math.Max(s.pointB.X, s.pointC.X)))
		pMax.Y = math.Max(pMax.Y, math.Max(s.pointA.Y, math.Max(s.pointB.Y, s.pointC.Y)))
		pMax.Z = math.Max(pMax.Z, math.Max(s.pointA.Z, math.Max(s.pointB.Z, s.pointC.Z)))

		if i == len(stlFile.Triangles)-1 || (idx+1)%boundingBoxMaxSize == 0 {
			shapes = append(shapes, &boundingBox{
				pMin:   pMin,
				pMax:   pMax,
				shapes: boundingBoxShapes,
			})
			boundingBoxShapes = make([]shape, 0, boundingBoxMaxSize)
			pMin = pMinDefault
			pMax = pMaxDefault
		}
	}

	fmt.Printf("Loaded stl file %s, %v triangles\n", stlFileName, len(shapes))
	return shapes
}
