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
	return genericShowCaseWithMirrorWalls(is, true, fromStlFile(
		"Istareyn/low-poly-stanford-bunny/Bunny-LowPoly.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: -45, Z: 0})
			(*sh).translate(r3.Vec{X: 0, Y: -1, Z: 0})
		},
	))
}

func bellsprout(is imageSpec) scene {
	return genericShowCaseWithMirrorWalls(is, true, fromStlFile(
		"Philin_theBlank/bellsprout-with-flower-pot/Bellsprout_in_Flower_Pot.stl",
		func(sh *shape) {
			(*sh).scale(0.25)
			(*sh).rotate(r3.Vec{X: -90, Y: 165, Z: 0})
			(*sh).translate(r3.Vec{X: 0.20, Y: -0.5, Z: -1})
		},
	))
}

func koala(is imageSpec) scene {
	return cornerGenericShowCase(is, fromStlFile(
		"TroySlatton/lyman-from-animal-crossing/Lyman.stl",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: 0, Y: 210, Z: 0})
			(*sh).translate(r3.Vec{X: 0.5, Y: -1, Z: -0.35})
		},
	))
}

func tesla(is imageSpec) scene {
	return genericShowCaseWithMirrorWalls(is, true, fromStlFile(
		"Sim3D_/tesla-model-3-for-3d-printing/solid/Tesla Model 3.STL",
		func(sh *shape) {
			(*sh).scale(0.015)
			(*sh).rotate(r3.Vec{X: -90, Y: 240, Z: 0})
			(*sh).translate(r3.Vec{X: 0.20, Y: -1, Z: 0})
		},
	))
}

func genericShowCaseWithMirrorWalls(is imageSpec, withWalls bool, centerShapes []shape) scene {
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
	shapes := append(centerShapes, floorRoof(-1, 3, phongBlinn{specValue: 0.0, specShininess: 0.0, color: r3.Vec{X: 255.0 / 255.0, Y: 235.0 / 255.0, Z: 205.0 / 255.0}})...)
	if withWalls {
		shapes = append(shapes, walls(3, metal{albedo: r3.Vec{X: 0.75, Y: 0.75, Z: 0.75}}, true, true, true, true)...)
	}
	lights := []light{
		spotLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 0 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: -1,
				Y: 2.9,
				Z: -1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: -1, Y: 2.9, Z: -1}),
			angle:     25,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 255 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 1,
				Y: 2.9,
				Z: -1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: 1, Y: 2.9, Z: -1}),
			angle:     25,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 0 / 255.0,
				Z: 255 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 0,
				Y: 2.9,
				Z: 1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: 0, Y: 2.9, Z: 1}),
			angle:     25,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 0 / 255.0,
				Z: 255 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: -1,
				Y: -1,
				Z: 1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: -1, Y: -1, Z: 1}),
			angle:     25,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 0 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 1,
				Y: -1,
				Z: 1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: 1, Y: -1, Z: 1}),
			angle:     25,
		},
		spotLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 255 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 0,
				Y: -1,
				Z: -1,
			},
			direction: r3.Sub(r3.Vec{X: 0, Y: 0, Z: 0}, r3.Vec{X: 0, Y: -1, Z: -1}),
			angle:     25,
		},
	}
	return scene{
		camera: &cam,
		shapes: &shapes,
		lights: &lights,
	}
}

func cornerGenericShowCase(is imageSpec, centerShapes []shape) scene {
	radius := 3.0
	lightDist := radius / 2
	lookFrom := r3.Vec{X: -(radius - 1), Y: 0, Z: -(radius - 1)}
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
	shapes := append(centerShapes, floorRoof(-1, radius, phongBlinn{specValue: 1.0, specShininess: 0, color: r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}})...)
	shapes = append(shapes, walls(radius, phongBlinn{specValue: 1.0, specShininess: 0, color: r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}}, true, true, true, true)...)
	lights := []light{
		pointLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 255 / 255.0,
				Z: 255 / 255.0,
			},
			lightIntensity: 2.0,
			position: r3.Vec{
				X: 0,
				Y: lightDist,
				Z: 0,
			},
		},
		pointLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 0 / 255.0,
				Z: 255 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: -lightDist,
				Y: 0,
				Z: -lightDist,
			},
		},
		pointLight{
			colorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 0 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: lightDist,
				Y: 0,
				Z: -lightDist,
			},
		},
		pointLight{
			colorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 255 / 255.0,
				Z: 0 / 255.0,
			},
			lightIntensity: 1.0,
			position: r3.Vec{
				X: 0,
				Y: 0,
				Z: lightDist,
			},
		},
	}
	return scene{
		camera: &cam,
		shapes: &shapes,
		lights: &lights,
	}
}

func floorRoof(yCoordFloor, yCoordRoof float64, mat material) []shape {
	return []shape{
		&triangle{
			pointA:      r3.Vec{X: 50, Y: yCoordFloor, Z: 50},
			pointB:      r3.Vec{X: 50, Y: yCoordFloor, Z: -50},
			pointC:      r3.Vec{X: -50, Y: yCoordFloor, Z: 50},
			singleSided: true,
			mat:         mat,
		},
		&triangle{
			pointA:      r3.Vec{X: -50, Y: yCoordFloor, Z: -50},
			pointB:      r3.Vec{X: -50, Y: yCoordFloor, Z: 50},
			pointC:      r3.Vec{X: 50, Y: yCoordFloor, Z: -50},
			singleSided: true,
			mat:         mat,
		},
		&triangle{
			pointA:      r3.Vec{X: 50, Y: yCoordRoof, Z: 50},
			pointB:      r3.Vec{X: -50, Y: yCoordRoof, Z: 50},
			pointC:      r3.Vec{X: 50, Y: yCoordRoof, Z: -50},
			singleSided: true,
			mat:         mat,
		},
		&triangle{
			pointA:      r3.Vec{X: -50, Y: yCoordRoof, Z: -50},
			pointB:      r3.Vec{X: 50, Y: yCoordRoof, Z: -50},
			pointC:      r3.Vec{X: -50, Y: yCoordRoof, Z: 50},
			singleSided: true,
			mat:         mat,
		},
	}
}

func walls(radius float64, mat material, front, back, left, right bool) []shape {
	shapes := make([]shape, 0, 8)
	if front {
		shapes = append(shapes, []shape{
			&triangle{
				pointA:      r3.Vec{X: radius, Y: radius, Z: radius},
				pointB:      r3.Vec{X: radius, Y: -radius, Z: radius},
				pointC:      r3.Vec{X: -radius, Y: radius, Z: radius},
				singleSided: true,
				mat:         mat,
			},
			&triangle{
				pointA:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				pointB:      r3.Vec{X: -radius, Y: radius, Z: radius},
				pointC:      r3.Vec{X: radius, Y: -radius, Z: radius},
				singleSided: true,
				mat:         mat,
			},
		}...)
	}
	if back {
		shapes = append(shapes, []shape{
			&triangle{
				pointA:      r3.Vec{X: radius, Y: radius, Z: -radius},
				pointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				pointC:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				singleSided: true,
				mat:         mat,
			},
			&triangle{
				pointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
				pointB:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				pointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				singleSided: true,
				mat:         mat,
			},
		}...)
	}
	if left {
		shapes = append(shapes, []shape{
			&triangle{
				pointA:      r3.Vec{X: radius, Y: radius, Z: radius},
				pointB:      r3.Vec{X: radius, Y: radius, Z: -radius},
				pointC:      r3.Vec{X: radius, Y: -radius, Z: radius},
				singleSided: true,
				mat:         mat,
			},
			&triangle{
				pointA:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				pointB:      r3.Vec{X: radius, Y: -radius, Z: radius},
				pointC:      r3.Vec{X: radius, Y: radius, Z: -radius},
				singleSided: true,
				mat:         mat,
			},
		}...)
	}
	if right {
		shapes = append(shapes, []shape{
			&triangle{
				pointA:      r3.Vec{X: -radius, Y: radius, Z: radius},
				pointB:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				pointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				singleSided: true,
				mat:         mat,
			},
			&triangle{
				pointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
				pointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				pointC:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				singleSided: true,
				mat:         mat,
			},
		}...)
	}
	return shapes
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
		X: math.MaxFloat64,
		Y: math.MaxFloat64,
		Z: math.MaxFloat64,
	}
	pMaxDefault := r3.Vec{
		X: math.MinInt64,
		Y: math.MinInt64,
		Z: math.MinInt64,
	}
	pMin := pMinDefault
	pMax := pMaxDefault
	shapes := make([]shape, 0, len(stlFile.Triangles))
	for i, stlTriangle := range stlFile.Triangles {
		s := triangle{
			pointA:      r3.Vec{X: float64(stlTriangle.Vertices[0][0]), Y: float64(stlTriangle.Vertices[0][1]), Z: float64(stlTriangle.Vertices[0][2])},
			pointB:      r3.Vec{X: float64(stlTriangle.Vertices[1][0]), Y: float64(stlTriangle.Vertices[1][1]), Z: float64(stlTriangle.Vertices[1][2])},
			pointC:      r3.Vec{X: float64(stlTriangle.Vertices[2][0]), Y: float64(stlTriangle.Vertices[2][1]), Z: float64(stlTriangle.Vertices[2][2])},
			singleSided: false,
			mat: dielectric{
				refractiveIndex: 0,
			},
			//mat: phongBlinn{
			//	specValue:     2,
			//	specShininess: 5,
			//	color:         r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()},
			//},
		}
		shapes = append(shapes, &s)
		mutator(&shapes[i])

		// calculate bounding box
		pMin.X = math.Min(pMin.X, math.Min(s.pointA.X, math.Min(s.pointB.X, s.pointC.X)))
		pMin.Y = math.Min(pMin.Y, math.Min(s.pointA.Y, math.Min(s.pointB.Y, s.pointC.Y)))
		pMin.Z = math.Min(pMin.Z, math.Min(s.pointA.Z, math.Min(s.pointB.Z, s.pointC.Z)))
		pMax.X = math.Max(pMax.X, math.Max(s.pointA.X, math.Max(s.pointB.X, s.pointC.X)))
		pMax.Y = math.Max(pMax.Y, math.Max(s.pointA.Y, math.Max(s.pointB.Y, s.pointC.Y)))
		pMax.Z = math.Max(pMax.Z, math.Max(s.pointA.Z, math.Max(s.pointB.Z, s.pointC.Z)))
	}

	fmt.Printf("Loaded stl file %s, %v triangles\n", stlFileName, len(shapes))
	return shapes
}
