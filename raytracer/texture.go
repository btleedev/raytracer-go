package raytracer

import (
	"fmt"
	"gonum.org/v1/gonum/spatial/r3"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

type texture interface {
	getColorFrac(u, v float64) r3.Vec
}

type CheckersTexture struct {
	ColorFrac1     r3.Vec
	ColorFrac2     r3.Vec
	CheckersWidth  float64
	CheckersHeight float64
}

type ImageTexture struct {
	Img *image.RGBA
}

func (t CheckersTexture) getColorFrac(u, v float64) r3.Vec {
	u2 := int(math.Floor(u * t.CheckersWidth))
	v2 := int(math.Floor(v * t.CheckersHeight))

	if (u2+v2)%2 == 0 {
		return t.ColorFrac1
	} else {
		return t.ColorFrac2
	}
}

func (i ImageTexture) getColorFrac(u, v float64) r3.Vec {
	u2 := int(math.Floor(u * float64(i.Img.Bounds().Size().X)))
	v2 := int(math.Floor(v * float64(i.Img.Bounds().Size().Y)))
	r, g, b, a := i.Img.At(u2, v2).RGBA()
	// numbers from [0, 65535], make it to [0, 255]
	r256 := float64(r) / 255.99
	g256 := float64(g) / 255.99
	b256 := float64(b) / 255.99
	a1 := math.Min(1.0, float64(a)/(255.99*255.99))

	return r3.Vec{
		X: (1.0-a1)*backgroundColorFracR + (a1*r256)/255.99,
		Y: (1.0-a1)*backgroundColorFracG + (a1*g256)/255.99,
		Z: (1.0-a1)*backgroundColorFracB + (a1*b256)/255.99,
	}
}

func LoadRGBAImage(file io.Reader) (*image.RGBA, error) {
	img, _, err := image.Decode(file)

	fmt.Printf("%v\n", err)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	loadedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixelId := ((y * width) + x) * 4
			r, g, b, a := img.At(x, y).RGBA()
			// not sure why above is uint32, and pixels are uint8, so do conversion
			r256 := uint8(math.Floor(math.Min(255, float64(r)/255.99)))
			g256 := uint8(math.Floor(math.Min(255, float64(g)/255.99)))
			b256 := uint8(math.Floor(math.Min(255, float64(b)/255.99)))
			a256 := uint8(math.Floor(math.Min(255, float64(a)/255.99)))
			loadedImg.Pix[pixelId+0] = r256
			loadedImg.Pix[pixelId+1] = g256
			loadedImg.Pix[pixelId+2] = b256
			loadedImg.Pix[pixelId+3] = a256
		}
	}

	return loadedImg, nil
}
