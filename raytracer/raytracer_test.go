package raytracer

import (
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"
	"testing"
	"time"
)

// simple regression that contains all shapes, materials, and lights and compares result image to an expected one
func TestRegression(t *testing.T) {
	is, sc, exp := exampleRegression640x380(t)
	is.BvhTraversalAlgorithm = DepthFirstSearch
	imgDfs, imgDfsDuration := renderImage(is, sc)
	is.BvhTraversalAlgorithm = Dijkstra
	imgDjikstras, imgDjikstrasDuration := renderImage(is, sc)

	fmt.Println()
	fmt.Printf("Djikstras algorithm render time: %v\n", imgDjikstrasDuration.String())
	fmt.Printf("DepthFirstSearch algorithm render time: %v\n", imgDfsDuration.String())

	fmt.Printf("Comparing image rendered by Djikstras algorithm to expected image")
	compareImages(t, imgDjikstras, exp)
	fmt.Printf("Comparing image rendered by DepthFirstSearch algorithm to expected image")
	compareImages(t, imgDfs, exp)
	fmt.Printf("Comparing image rendered by Djikstras algorithm to image rendered by DepthFirstSearch algorithm")
	compareImages(t, imgDjikstras, imgDfs)
}

func compareImages(t *testing.T, img *image.RGBA, exp *image.RGBA) {
	width := exp.Rect.Max.X - exp.Rect.Min.X
	height := exp.Rect.Max.Y - exp.Rect.Min.Y
	// there are some random logic (eg anti-aliasing, di-electric material)
	// anti-aliasing should hopefully eliminate randomness, but we need to add an acceptable delta
	antiAliasingDelta := uint32(20 * 257)
	maximumDifferentPixelsAllowed := int(float64(width*height) * 0.01)

	if img.Rect.Min.X != exp.Rect.Min.X ||
		img.Rect.Min.Y != exp.Rect.Min.Y ||
		img.Rect.Max.X != exp.Rect.Max.X ||
		img.Rect.Max.Y != exp.Rect.Max.Y {

		t.Error("Generated image does not match expected image size")
	}

	differentPixels := 0
	for i := exp.Rect.Min.X; i <= exp.Rect.Max.X; i++ {
		for j := exp.Rect.Min.Y; j <= exp.Rect.Max.Y; j++ {
			imgRgba := img.At(i, j)
			expRgba := exp.At(i, j)
			ir, ig, ib, ia := imgRgba.RGBA()
			er, eg, eb, ea := expRgba.RGBA()
			if diffu32(ir, er) > antiAliasingDelta ||
				diffu32(ig, eg) > antiAliasingDelta ||
				diffu32(ib, eb) > antiAliasingDelta ||
				ia != ea {

				differentPixels++
				// t.Errorf("Pixel (%d, %d) are too different, was (%d, %d, %d, %d) but expected (%d, %d, %d, %d)", i, j, ir, ig, ib, ia, er, eg, eb, ea)
			}
		}
	}

	fmt.Printf(
		"Image was the same in %d pixels with %d delta, but was different in %d (%.2f%%) pixels\n",
		width*height-differentPixels,
		antiAliasingDelta,
		differentPixels,
		100.0*float64(differentPixels)/float64(width*height),
	)
	if differentPixels >= maximumDifferentPixelsAllowed {
		t.Errorf("Different pixels exceeded threshold of %d, (%.2f%%)",
			maximumDifferentPixelsAllowed,
			100.0*float64(maximumDifferentPixelsAllowed)/float64(width*height),
		)
	}
	fmt.Println()
}

func renderImage(imageSpec ImageSpec, scene Scene) (i *image.RGBA, d time.Duration) {
	startTime := time.Now()
	img := GenerateImage(imageSpec, scene)
	return img, time.Since(startTime)
}

func exampleRegression640x380(t *testing.T) (is ImageSpec, sc Scene, exp *image.RGBA) {
	imageSpec, scene := ExampleRegression(640, 380)
	fileName := "../samples_images/code_example.png"
	file, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	expectedImage, err := loadRGBAImage(file)
	if err != nil {
		t.Error(err)
	}
	return imageSpec, scene, expectedImage
}

func loadRGBAImage(file io.Reader) (*image.RGBA, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	loadedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixelId := ((y * width) + x) * 4
			r, g, b, a := img.At(x, y).RGBA()
			loadedImg.Pix[pixelId+0] = uint8(r)
			loadedImg.Pix[pixelId+1] = uint8(g)
			loadedImg.Pix[pixelId+2] = uint8(b)
			loadedImg.Pix[pixelId+3] = uint8(a)
		}
	}

	return loadedImg, nil
}

func diffu32(i, j uint32) uint32 {
	if i > j {
		return i - j
	}
	return j - i
}
