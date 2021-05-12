package raytracer

import "gonum.org/v1/gonum/spatial/r3"

func ExampleRegression(width, height int) (is ImageSpec, sc Scene) {
	floorRadius := 100.0
	centerPiecesRadius := 2.0
	backMirrorRadius := 4 * centerPiecesRadius
	backMirrorBorder := centerPiecesRadius / 2

	cameraLookFrom := r3.Vec{X: 0, Y: 3 * centerPiecesRadius, Z: -8}
	cameraLookAt := r3.Vec{X: 0, Y: 1.5 * centerPiecesRadius, Z: 0}
	cameraUp := r3.Vec{X: 0, Y: 1, Z: 0}
	cameraFocusPoint := cameraLookAt
	cameraAperature := 0.015
	cameraFovDegrees := 60.0

	shapes := []Shape{
		// centerpieces
		&Sphere{
			Center: r3.Vec{X: 2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: Dielectric{
				RefractiveIndex: 1.0,
			},
		},
		&Sphere{
			Center: r3.Vec{X: 0, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: PhongBlinn{
				SpecHardness:  1,
				SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1},
				Color:         r3.Vec{X: 1, Y: 1, Z: 1},
			},
		},
		&Sphere{
			Center: r3.Vec{X: -2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: Dielectric{
				RefractiveIndex: 0,
			},
		},

		// floor
		&TrianglePlane{
			PointA:      r3.Vec{X: -floorRadius, Y: 0, Z: -floorRadius},
			PointB:      r3.Vec{X: -floorRadius, Y: 0, Z: floorRadius},
			PointC:      r3.Vec{X: floorRadius, Y: 0, Z: -floorRadius},
			SingleSided: true,
			Mat: PhongBlinn{
				Color:         r3.Vec{X: 0, Y: 0, Z: 0},
				SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1},
				SpecHardness:  1,
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: floorRadius, Y: 0, Z: floorRadius},
			PointB:      r3.Vec{X: floorRadius, Y: 0, Z: -floorRadius},
			PointC:      r3.Vec{X: -floorRadius, Y: 0, Z: floorRadius},
			SingleSided: true,
			Mat: PhongBlinn{
				Color:         r3.Vec{X: 0, Y: 0, Z: 0},
				SpecularColor: r3.Vec{X: 1, Y: 1, Z: 1},
				SpecHardness:  1,
			},
		},

		// back mirror
		&TrianglePlane{
			PointA:      r3.Vec{X: backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			PointB:      r3.Vec{X: backMirrorRadius, Y: 0, Z: backMirrorRadius},
			PointC:      r3.Vec{X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			SingleSided: true,
			Mat: Standard{
				Color: r3.Vec{X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0},
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: -backMirrorRadius, Y: 0, Z: backMirrorRadius},
			PointB:      r3.Vec{X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			PointC:      r3.Vec{X: backMirrorRadius, Y: 0, Z: backMirrorRadius},
			SingleSided: true,
			Mat: Standard{
				Color: r3.Vec{X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0},
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: backMirrorRadius - backMirrorBorder, Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			PointB:      r3.Vec{X: backMirrorRadius - backMirrorBorder, Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			PointC:      r3.Vec{X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			SingleSided: true,
			Mat: Metal{
				Albedo: r3.Vec{X: 1, Y: 1, Z: 1},
				Fuzz:   0,
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			PointB:      r3.Vec{X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			PointC:      r3.Vec{X: backMirrorRadius - backMirrorBorder, Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
			SingleSided: true,
			Mat: Metal{
				Albedo: r3.Vec{X: 1, Y: 1, Z: 1},
				Fuzz:   0,
			},
		},
	}
	lights := []Light{
		AmbientLight{
			ColorFrac: r3.Vec{
				X: 255 / 255.0,
				Y: 0 / 255.0,
				Z: 0 / 255.0,
			},
			LightIntensity: 0.2,
		},
		SpotLight{
			ColorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 255 / 255.0,
				Z: 0 / 255.0,
			},
			LightIntensity:         100,
			SpecularLightIntensity: 100,
			Position: r3.Vec{
				X: 6 * centerPiecesRadius,
				Y: 5 * centerPiecesRadius,
				Z: 3 * centerPiecesRadius,
			},
			Direction: r3.Vec{
				X: -6,
				Y: -5,
				Z: -3,
			},
			Angle: 30,
		},
		PointLight{
			ColorFrac: r3.Vec{
				X: 0 / 255.0,
				Y: 0 / 255.0,
				Z: 255 / 255.0,
			},
			LightIntensity:         100,
			SpecularLightIntensity: 100,
			Position: r3.Vec{
				X: -6 * centerPiecesRadius,
				Y: 5 * centerPiecesRadius,
				Z: -3 * centerPiecesRadius,
			},
		},
	}
	imageSpec := ImageSpec{
		Width:                           width,
		Height:                          height,
		AntiAliasingFactor:              32,
		RayTracingMaxDepth:              16,
		SoftShadowMonteCarloRepetitions: 16,
		WorkerCount:                     16,
	}
	scene := Scene{
		CameraLookFrom:   cameraLookFrom,
		CameraLookAt:     cameraLookAt,
		CameraUp:         cameraUp,
		CameraFocusPoint: cameraFocusPoint,
		CameraAperature:  cameraAperature,
		CameraFov:        cameraFovDegrees,
		Shapes:           shapes,
		Lights:           lights,
	}
	return imageSpec, scene
}
