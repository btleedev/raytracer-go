package main

import (
	"example.com/hello/raytracer"
	"fmt"
	"github.com/hschendel/stl"
	"gonum.org/v1/gonum/spatial/r3"
	"math"
	"math/rand"
	"sort"
)

func sample() (cameraLookFrom r3.Vec, cameraLookAt r3.Vec, cameraUp r3.Vec, cameraFocusPoint r3.Vec, s []raytracer.Shape, l []raytracer.Light) {
	radius := 20.0
	lookFrom := r3.Vec{X: 0, Y: 0, Z: -3 * radius}
	lookAt := r3.Vec{X: 0, Y: 0, Z: 0}
	up := r3.Vec{X: 0, Y: 1, Z: 0}
	shapes := []raytracer.Shape{
		&raytracer.Sphere{
			Center: r3.Vec{X: 0, Y: -radius + (radius / 3), Z: 0},
			Radius: radius / 3,
			Mat: raytracer.PhongBlinn{
				SpecHardness:  1,
				SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1},
				Color:         r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()},
			},
		},
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: 2 * radius / 3, Y: -radius + 0.01, Z: 2 * radius / 3},
			PointB:      r3.Vec{X: 2 * radius / 3, Y: -radius + 0.01, Z: -2 * radius / 3},
			PointC:      r3.Vec{X: -2 * radius / 3, Y: -radius + 0.01, Z: -2 * radius / 3},
			SingleSided: true,
			Mat: raytracer.Metal{
				Albedo: r3.Vec{X: 1.0, Y: 1.0, Z: 1.0},
				Fuzz:   0,
			},
		},
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: 2 * radius / 3, Y: -radius + 0.01, Z: 2 * radius / 3},
			PointB:      r3.Vec{X: -2 * radius / 3, Y: -radius + 0.01, Z: -2 * radius / 3},
			PointC:      r3.Vec{X: -2 * radius / 3, Y: -radius + 0.01, Z: 2 * radius / 3},
			SingleSided: true,
			Mat: raytracer.Metal{
				Albedo: r3.Vec{X: 1.0, Y: 1.0, Z: 1.0},
				Fuzz:   0,
			},
		},
	}
	shapes = append(shapes, floorRoof(-radius, radius, radius, raytracer.PhongBlinn{SpecHardness: 2, SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1}, Color: r3.Vec{X: rand.Float64(), Y: rand.Float64(), Z: rand.Float64()}})...)
	shapes = append(shapes, walls(radius, raytracer.Metal{Albedo: r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, Fuzz: 0}, true, false, false, false)...)
	shapes = append(shapes, walls(radius, raytracer.PhongBlinn{SpecHardness: 1, SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1}, Color: r3.Vec{X: 1, Y: 0, Z: 0}}, false, true, false, false)...)
	shapes = append(shapes, walls(radius, raytracer.PhongBlinn{SpecHardness: 1, SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1}, Color: r3.Vec{X: 0, Y: 1, Z: 0}}, false, false, true, false)...)
	shapes = append(shapes, walls(radius, raytracer.PhongBlinn{SpecHardness: 1, SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1}, Color: r3.Vec{X: 0, Y: 0, Z: 1}}, false, false, false, true)...)
	lights := []raytracer.Light{
		raytracer.AmbientLight{
			ColorFrac: r3.Vec{
				X: 1,
				Y: 1,
				Z: 1,
			},
			LightIntensity: 0.5,
		},
		raytracer.PointLight{
			ColorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 255 / 255.0,
				Z: 255 / 255.0,
			},
			LightIntensity:         100,
			SpecularLightIntensity: 100,
			Position: r3.Vec{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
	}
	return lookFrom, lookAt, up, lookAt, shapes, lights
}

func floorRoof(yCoordFloor, yCoordRoof, radius float64, mat raytracer.Material) []raytracer.Shape {
	return []raytracer.Shape{
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: radius, Y: yCoordFloor, Z: radius},
			PointB:      r3.Vec{X: radius, Y: yCoordFloor, Z: -radius},
			PointC:      r3.Vec{X: -radius, Y: yCoordFloor, Z: radius},
			SingleSided: true,
			Mat:         mat,
		},
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: -radius, Y: yCoordFloor, Z: -radius},
			PointB:      r3.Vec{X: -radius, Y: yCoordFloor, Z: radius},
			PointC:      r3.Vec{X: radius, Y: yCoordFloor, Z: -radius},
			SingleSided: true,
			Mat:         mat,
		},
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: radius, Y: yCoordRoof, Z: radius},
			PointB:      r3.Vec{X: -radius, Y: yCoordRoof, Z: radius},
			PointC:      r3.Vec{X: radius, Y: yCoordRoof, Z: -radius},
			SingleSided: true,
			Mat:         mat,
		},
		&raytracer.TrianglePlane{
			PointA:      r3.Vec{X: -radius, Y: yCoordRoof, Z: -radius},
			PointB:      r3.Vec{X: radius, Y: yCoordRoof, Z: -radius},
			PointC:      r3.Vec{X: -radius, Y: yCoordRoof, Z: radius},
			SingleSided: true,
			Mat:         mat,
		},
	}
}

func walls(radius float64, mat raytracer.Material, front, back, left, right bool) []raytracer.Shape {
	shapes := make([]raytracer.Shape, 0, 8)
	if front {
		shapes = append(shapes, []raytracer.Shape{
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: radius, Y: radius, Z: radius},
				PointB:      r3.Vec{X: radius, Y: -radius, Z: radius},
				PointC:      r3.Vec{X: -radius, Y: radius, Z: radius},
				SingleSided: true,
				Mat:         mat,
			},
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				PointB:      r3.Vec{X: -radius, Y: radius, Z: radius},
				PointC:      r3.Vec{X: radius, Y: -radius, Z: radius},
				SingleSided: true,
				Mat:         mat,
			},
		}...)
	}
	if back {
		shapes = append(shapes, []raytracer.Shape{
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: radius, Y: radius, Z: -radius},
				PointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				PointC:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				SingleSided: true,
				Mat:         mat,
			},
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
				PointB:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				PointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				SingleSided: true,
				Mat:         mat,
			},
		}...)
	}
	if left {
		shapes = append(shapes, []raytracer.Shape{
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: radius, Y: radius, Z: radius},
				PointB:      r3.Vec{X: radius, Y: radius, Z: -radius},
				PointC:      r3.Vec{X: radius, Y: -radius, Z: radius},
				SingleSided: true,
				Mat:         mat,
			},
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: radius, Y: -radius, Z: -radius},
				PointB:      r3.Vec{X: radius, Y: -radius, Z: radius},
				PointC:      r3.Vec{X: radius, Y: radius, Z: -radius},
				SingleSided: true,
				Mat:         mat,
			},
		}...)
	}
	if right {
		shapes = append(shapes, []raytracer.Shape{
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: -radius, Y: radius, Z: radius},
				PointB:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				PointC:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				SingleSided: true,
				Mat:         mat,
			},
			&raytracer.TrianglePlane{
				PointA:      r3.Vec{X: -radius, Y: -radius, Z: -radius},
				PointB:      r3.Vec{X: -radius, Y: radius, Z: -radius},
				PointC:      r3.Vec{X: -radius, Y: -radius, Z: radius},
				SingleSided: true,
				Mat:         mat,
			},
		}...)
	}
	return shapes
}

func fromStlFile(stlFileName string, mutator func(shape *raytracer.Shape)) []raytracer.Shape {
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
	shapes := make([]raytracer.Shape, 0, len(stlFile.Triangles))
	for i, stlTriangle := range stlFile.Triangles {
		s := raytracer.TrianglePlane{
			PointA:      r3.Vec{X: float64(stlTriangle.Vertices[0][0]), Y: float64(stlTriangle.Vertices[0][1]), Z: float64(stlTriangle.Vertices[0][2])},
			PointB:      r3.Vec{X: float64(stlTriangle.Vertices[1][0]), Y: float64(stlTriangle.Vertices[1][1]), Z: float64(stlTriangle.Vertices[1][2])},
			PointC:      r3.Vec{X: float64(stlTriangle.Vertices[2][0]), Y: float64(stlTriangle.Vertices[2][1]), Z: float64(stlTriangle.Vertices[2][2])},
			SingleSided: true,
			Mat: raytracer.Dielectric{
				RefractiveIndex: 0,
			},
		}
		shapes = append(shapes, &s)
		mutator(&shapes[i])

		// calculate bounding box
		pMin.X = math.Min(pMin.X, math.Min(s.PointA.X, math.Min(s.PointB.X, s.PointC.X)))
		pMin.Y = math.Min(pMin.Y, math.Min(s.PointA.Y, math.Min(s.PointB.Y, s.PointC.Y)))
		pMin.Z = math.Min(pMin.Z, math.Min(s.PointA.Z, math.Min(s.PointB.Z, s.PointC.Z)))
		pMax.X = math.Max(pMax.X, math.Max(s.PointA.X, math.Max(s.PointB.X, s.PointC.X)))
		pMax.Y = math.Max(pMax.Y, math.Max(s.PointA.Y, math.Max(s.PointB.Y, s.PointC.Y)))
		pMax.Z = math.Max(pMax.Z, math.Max(s.PointA.Z, math.Max(s.PointB.Z, s.PointC.Z)))
	}

	fmt.Printf("Loaded stl file %s, %v triangles\n", stlFileName, len(shapes))
	return shapes
}
