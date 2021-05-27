Ray tracer written in golang.

![bunny](samples_images/bunny.png "Stanford Bunny")
![lyman](samples_images/lyman.png "Lyman from Animal Crossing")

# Usage

See [Code Example](#code-example) for sample code.

```shell
go build -o raytracer-go
./raytracer-go

open ./out.png
```

![Code Example](samples_images/code_example.png "Code Example")

# Shapes

* Sphere
* Triangle plane

# Lighting

* Ambient
* Point
* Spot

# Materials

* Standard
* Metal
* Dielectric
* Phong-Blinn

# Features

* Acceleration structures (bounding volume hierarchy)
* Anti-Aliasing
* Camera FOV
* Camera Lens blur (aperature)
* Inverse square law decay for non-ambient lights
* Soft Shadows (Monte Carlo)
* Transformations (translate, scale, rotate)

# STL Models used in samples

* [LOW POLY STANFORD BUNNY](https://cults3d.com/en/3d-model/art/low-poly-stanford-bunny) - Author: Istareyn - LICENSE: https://creativecommons.org/licenses/by-nc/3.0/ - no changes to STL file were made
* [LYMAN FROM ANIMAL CROSSING](https://cults3d.com/en/3d-model/art/lyman-from-animal-crossing) - Author: Philin_theBlank - LICENSE: https://creativecommons.org/licenses/by-nc/3.0/ - no changes to STL file were made

# Code Example

```
    imageLocation := "out.png"

    floorRadius := 100.0
    centerPiecesRadius := 2.0
    backMirrorRadius := 4 * centerPiecesRadius
    backMirrorBorder := centerPiecesRadius / 2

    cameraLookFrom := r3.Vec{ X: 0, Y: 3 * centerPiecesRadius, Z: -8 }
    cameraLookAt := r3.Vec{ X: 0, Y: 1.5 * centerPiecesRadius, Z: 0 }
    cameraUp := r3.Vec{ X: 0, Y: 1, Z: 0 }
    cameraFocusPoint := cameraLookAt
    cameraAperature := 0.015
    cameraFovDegrees := 60.0

    shapes := []raytracer.Shape{
        // centerpieces
        &raytracer.Sphere{
            Center: r3.Vec{ X: 2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0 },
            Radius: centerPiecesRadius,
            Mat: raytracer.Dielectric{
                RefractiveIndex: 1.52,
            },
        },
        &raytracer.Sphere{
            Center: r3.Vec{ X: 0, Y: centerPiecesRadius, Z: 0 },
            Radius: centerPiecesRadius,
            Mat: raytracer.PhongBlinn{
                SpecHardness:      1,
                SpecularColorFrac: r3.Vec{ X: 1, Y: 1, Z: 1 },
                ColorFrac:         r3.Vec{ X: 1, Y: 1, Z: 1 },
            },
        },
        &raytracer.Sphere{
            Center: r3.Vec{ X: -2 * centerPiecesRadius, Y: centerPiecesRadius, Z: 0 },
            Radius: centerPiecesRadius,
            Mat: raytracer.Metal{
                Albedo: r3.Vec{ X: 1.0, Y: 1.0, Z: 1.0},
                Fuzz:   0,
            },
        },

        // floor
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: -floorRadius, Y: 0, Z: -floorRadius },
            PointB:      r3.Vec{ X: -floorRadius, Y: 0, Z: floorRadius },
            PointC:      r3.Vec{ X: floorRadius, Y: 0, Z: -floorRadius },
            SingleSided: true,
            Mat: raytracer.PhongBlinn{
                ColorFrac:         r3.Vec{ X: 0, Y: 0, Z: 0 },
                SpecularColorFrac: r3.Vec{ X: 1, Y: 1, Z: 1 },
                SpecHardness:      1,
            },
        },
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: floorRadius, Y: 0, Z: floorRadius },
            PointB:      r3.Vec{ X: floorRadius, Y: 0, Z: -floorRadius },
            PointC:      r3.Vec{ X: -floorRadius, Y: 0, Z: floorRadius },
            SingleSided: true,
            Mat: raytracer.PhongBlinn{
                ColorFrac:         r3.Vec{ X: 0, Y: 0, Z: 0 },
                SpecularColorFrac: r3.Vec{ X: 1, Y: 1, Z: 1 },
                SpecHardness:      1,
            },
        },

        // back mirror
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius },
            PointB:      r3.Vec{ X: backMirrorRadius, Y: 0, Z: backMirrorRadius },
            PointC:      r3.Vec{ X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius },
            SingleSided: true,
            Mat: raytracer.Standard{
                ColorFrac: r3.Vec{ X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0 },
            },
        },
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: -backMirrorRadius, Y: 0, Z: backMirrorRadius },
            PointB:      r3.Vec{ X: -backMirrorRadius, Y: backMirrorRadius, Z: backMirrorRadius },
            PointC:      r3.Vec{ X: backMirrorRadius, Y: 0, Z: backMirrorRadius },
            SingleSided: true,
            Mat: raytracer.Standard{
                ColorFrac: r3.Vec{ X: 150 / 255.0, Y: 111 / 255.0, Z: 51 / 255.0 },
            },
        },
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: backMirrorRadius - backMirrorBorder, Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            PointB:      r3.Vec{ X: backMirrorRadius - backMirrorBorder, Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            PointC:      r3.Vec{ X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            SingleSided: true,
            Mat: raytracer.Metal{
                Albedo: r3.Vec{ X: 1, Y: 1, Z: 1 },
                Fuzz:   0,
            },
        },
        &raytracer.TrianglePlane{
            PointA:      r3.Vec{ X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            PointB:      r3.Vec{ X: -(backMirrorRadius - backMirrorBorder), Y: backMirrorRadius - backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            PointC:      r3.Vec{ X: backMirrorRadius - backMirrorBorder, Y: backMirrorBorder, Z: backMirrorRadius - backMirrorBorder},
            SingleSided: true,
            Mat: raytracer.Metal{
                Albedo: r3.Vec{ X: 1, Y: 1, Z: 1 },
                Fuzz:   0,
            },
        },
    }
    lights := []raytracer.Light{
        raytracer.AmbientLight{
            ColorFrac: r3.Vec{
                X: 255 / 255.0,
                Y: 0 / 255.0,
                Z: 0 / 255.0,
            },
            LightIntensity: 0.2,
        },
        raytracer.SpotLight{
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
            LookAt: r3.Vec{ 
                X: 0, 
                Y: 0, 
                Z: 0,
            },
            Angle: 30,
            InverseSquareLawDecayFactor: 1.0,
        },
        raytracer.PointLight{
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
            InverseSquareLawDecayFactor: 1.0,
        },
    }
    imageSpec := raytracer.ImageSpec{
        Width:                           640,
        Height:                          380,
        AntiAliasingFactor:              32,
        RayTracingMaxDepth:              16,
        SoftShadowMonteCarloRepetitions: 16,
        WorkerCount:                     16,
        BvhTraversalAlgorithm:           raytracer.Dijkstra,
    }
    scene := raytracer.Scene{
        CameraLookFrom:                  cameraLookFrom,
        CameraLookAt:                    cameraLookAt,
        CameraUp:                        cameraUp,
        CameraFocusPoint:                cameraFocusPoint,
        CameraAperature:                 cameraAperature,
        CameraFov:                       cameraFovDegrees,
        Shapes:                          shapes,
        Lights:                          lights,
    }
    myImage := raytracer.GenerateImage(imageSpec, scene)

    outputFile, err := os.Create(imageLocation)
    if err != nil {
        panic("failed to create image")
    }
    defer outputFile.Close()
    png.Encode(outputFile, myImage)
```