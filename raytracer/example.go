package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"os"
)

func ExampleRegression(width, height int, repoBaseDir string) (is ImageSpec, sc Scene) {
	floorRadius := 100.0
	centerPiecesRadius := 2.0
	backMirrorRadius := 4 * centerPiecesRadius
	backMirrorBorder := centerPiecesRadius / 2

	cameraLookFrom := r3.Vec{X: 0, Y: 3 * centerPiecesRadius, Z: -5}
	cameraLookAt := r3.Vec{X: 0, Y: 2 * centerPiecesRadius, Z: 0}
	cameraUp := r3.Vec{X: 0, Y: 1, Z: 0}
	cameraFocusPoint := cameraLookAt
	cameraAperature := 0.015
	cameraFovDegrees := 60.0

	texturePlane := CheckersTexture{
		ColorFrac1:     r3.Vec{X: 0, Y: 1, Z: 0},
		ColorFrac2:     r3.Vec{X: 0, Y: 0, Z: 1},
		CheckersWidth:  100.0,
		CheckersHeight: 100.0,
	}
	textureLeftSphere := CheckersTexture{
		ColorFrac1:     r3.Vec{X: 0, Y: 0, Z: 0},
		ColorFrac2:     r3.Vec{X: 1, Y: 1, Z: 1},
		CheckersWidth:  10.0,
		CheckersHeight: 10.0,
	}
	textureRightSphereFileName := fmt.Sprintf("%s/%s", repoBaseDir, "samples_textures/Tiles075_1K_Color.jpg")
	textureRightSphereFile, err := os.Open(textureRightSphereFileName)
	if err != nil {
		panic(err)
	}
	defer textureRightSphereFile.Close()

	textureRightSphereTexture, err := LoadRGBAImage(textureRightSphereFile)
	if err != nil {
		panic(err)
	}
	textureRightSphere := ImageTexture{
		Img: textureRightSphereTexture,
	}

	shapes := []Shape{
		// centerpieces
		&Sphere{
			Center: r3.Vec{X: 4 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: Standard{
				Texture: textureLeftSphere,
			},
		},
		&Sphere{
			Center: r3.Vec{X: 2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: Dielectric{
				RefractiveIndex: 1.52,
			},
		},
		&Sphere{
			Center: r3.Vec{X: 0, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: PhongBlinn{
				SpecHardness:      1,
				SpecularColorFrac: r3.Vec{X: 1, Y: 1, Z: 1},
				ColorFrac:         r3.Vec{X: 1, Y: 1, Z: 1},
			},
		},
		&Sphere{
			Center: r3.Vec{X: -2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: Metal{
				Albedo: r3.Vec{X: 1.0, Y: 1.0, Z: 1.0},
				Fuzz:   0,
			},
		},
		&Sphere{
			Center: r3.Vec{X: -4 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0},
			Radius: centerPiecesRadius,
			Mat: PhongBlinn{
				SpecHardness:      1,
				SpecularColorFrac: r3.Vec{X: 1, Y: 1, Z: 1},
				Texture:           textureRightSphere,
			},
		},

		// floor
		&TrianglePlane{
			PointA:      r3.Vec{X: -floorRadius, Y: 0, Z: -floorRadius},
			PointB:      r3.Vec{X: -floorRadius, Y: 0, Z: floorRadius},
			PointC:      r3.Vec{X: floorRadius, Y: 0, Z: -floorRadius},
			SingleSided: true,
			Mat: PhongBlinn{
				ColorFrac:         r3.Vec{X: 0, Y: 0, Z: 0},
				SpecularColorFrac: r3.Vec{X: 1, Y: 1, Z: 1},
				SpecHardness:      1,
				Texture:           texturePlane,
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: floorRadius, Y: 0, Z: floorRadius},
			PointB:      r3.Vec{X: floorRadius, Y: 0, Z: -floorRadius},
			PointC:      r3.Vec{X: -floorRadius, Y: 0, Z: floorRadius},
			SingleSided: true,
			Mat: PhongBlinn{
				ColorFrac:         r3.Vec{X: 0, Y: 0, Z: 0},
				SpecularColorFrac: r3.Vec{X: 1, Y: 1, Z: 1},
				SpecHardness:      1,
				Texture:           texturePlane,
			},
		},

		// back mirror
		&TrianglePlane{
			PointA:      r3.Vec{X: backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			PointB:      r3.Vec{X: backMirrorRadius, Y: 0, Z: backMirrorRadius},
			PointC:      r3.Vec{X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			SingleSided: true,
			Mat: Standard{
				ColorFrac: r3.Vec{X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0},
			},
		},
		&TrianglePlane{
			PointA:      r3.Vec{X: -backMirrorRadius, Y: 0, Z: backMirrorRadius},
			PointB:      r3.Vec{X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius},
			PointC:      r3.Vec{X: backMirrorRadius, Y: 0, Z: backMirrorRadius},
			SingleSided: true,
			Mat: Standard{
				ColorFrac: r3.Vec{X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0},
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
				X: 171 / 255.0,
				Y: 137 / 255.0,
				Z: 255 / 255.0,
			},
			LightIntensity:         100,
			SpecularLightIntensity: 100,
			Position: r3.Vec{
				X: 6 * centerPiecesRadius,
				Y: 5 * centerPiecesRadius,
				Z: -3 * centerPiecesRadius,
			},
			LookAt: r3.Vec{
				X: 0,
				Y: 0,
				Z: 0,
			},
			Angle:                       30,
			InverseSquareLawDecayFactor: 1.0,
		},
		PointLight{
			ColorFrac: r3.Vec{
				X: 67 / 255.0,
				Y: 163 / 255.0,
				Z: 241 / 255.0,
			},
			LightIntensity:         100,
			SpecularLightIntensity: 10,
			Position: r3.Vec{
				X: -4 * centerPiecesRadius,
				Y: centerPiecesRadius,
				Z: 3 * centerPiecesRadius,
			},
			InverseSquareLawDecayFactor: 0.5,
		},
	}
	imageSpec := ImageSpec{
		Width:                           width,
		Height:                          height,
		AntiAliasingFactor:              32,
		RayTracingMaxDepth:              16,
		SoftShadowMonteCarloRepetitions: 16,
		WorkerCount:                     16,
		BvhTraversalAlgorithm:           Dijkstra,
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
