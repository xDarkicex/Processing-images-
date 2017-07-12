package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	img, _ := loadImg("original.jpg")
	gray := rgbaToGray(img)
	dithered := FloydSteinbergDither(gray)
	// OGdithered := GridDither(gray, 2, 3, 9)

	// Save as gray.png
	f, _ := os.Create("dithered.png")
	defer f.Close()
	png.Encode(f, dithered)
}

func loadImg(filepath string) (image.Image, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()
	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

func blackOrWhite(g color.Gray) color.Gray {
	if g.Y < 127 {
		return color.Gray{0} //black
	}
	return color.Gray{255} //white
}

func newWhite(bounds image.Rectangle) *image.Gray {
	white := image.NewGray(bounds)
	for i := range white.Pix {
		white.Pix[i] = 255
	}
	return white
}

func ThresholdDither(gray *image.Gray) *image.Gray {
	var (
		bounds   = gray.Bounds()
		dithered = image.NewGray(bounds)
		width    = bounds.Dx()
		height   = bounds.Dy()
	)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			var c = blackOrWhite(gray.GrayAt(i, j))
			dithered.SetGray(i, j, c)
		}
	}
	return dithered
}

func avgIntensity(gray *image.Gray) float64 {
	var sum float64
	for _, pix := range gray.Pix {
		sum += float64(pix)
	}
	average := sum / float64(len(gray.Pix)*256)
	return (1 / (1 + (math.Pow(math.E, -10*(average-.5)))))
}

// GridDither Original dither function its alright
func GridDither(gray *image.Gray, cellsize int, alpha, gamma float64) *image.Gray {
	var (
		bounds   = gray.Bounds()
		dithered = newWhite(bounds)
		width    = bounds.Dx()
		height   = bounds.Dy()
		rng      = rand.New(rand.NewSource(time.Now().UnixNano()))
	)
	go dither(0, width, height, cellsize, gray, dithered, gamma, alpha, rng)
	dither(cellsize, width, height, cellsize, gray, dithered, gamma, alpha, rng)
	return dithered
}

// Function to thread Dither function
func dither(offset, width, height, cellsize int, gray, dithered *image.Gray, gamma, alpha float64, rng *rand.Rand) {
	for i := offset; i < width; i += cellsize * 2 {
		for j := offset; j < height; j += cellsize * 2 {
			var (
				cell = rgbaToGray(gray.SubImage(image.Rect(i, j, i+cellsize, j+cellsize)))
				mu   = avgIntensity(cell)
				n    = math.Pow((1-mu)*gamma, 2) / 3
			)
			if n < alpha {
				n = 0
			}
			for k := 0; k < int(n); k++ {
				var (
					x = randInt(i, min(i+k, width), rng)
					y = randInt(j, min(j+k, height), rng)
				)
				dithered.SetGray(x, y, color.Gray{0})
			}
		}
	}
}
func randInt(min, max int, rng *rand.Rand) int {
	if max-min == 0 {
		return 0
	}
	return rng.Intn(max-min) + min
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

//FloydSteinbergDither Fancy ass dither function.
func FloydSteinbergDither(gray *image.Gray) *image.Gray {
	var (
		bounds   = gray.Bounds()
		width    = bounds.Dx()
		height   = bounds.Dy()
		dithered = copyGray(gray)
	)
	// go fDither(0, dithered, gray, width, height)
	fDither(0, dithered, gray, width, height)
	return dithered
}

func copyGray(gray *image.Gray) *image.Gray {
	var clone = image.NewGray(gray.Bounds())
	copy(clone.Pix, gray.Pix)
	return clone
}

func i16ToUI8(x int16) uint8 {
	switch {
	case x < 1:
		return uint8(0)
	case x > 254:
		return uint8(255)
	}
	return uint8(x)
}

func fDither(offset int, dithered, gray *image.Gray, width, height int) {
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			var oldpixel = dithered.GrayAt(i, j)
			var newpixel = blackOrWhite(oldpixel)
			dithered.SetGray(i, j, newpixel)
			var quant = (int16(oldpixel.Y) - int16(newpixel.Y)) / 16
			dithered.SetGray(i+1, j, color.Gray{i16ToUI8(int16(dithered.GrayAt(i+1, j).Y) + 7*quant)})
			dithered.SetGray(i-1, j+1, color.Gray{i16ToUI8(int16(dithered.GrayAt(i-1, j+1).Y) + 3*quant)})
			dithered.SetGray(i, j+1, color.Gray{i16ToUI8(int16(dithered.GrayAt(i, j+1).Y) + 5*quant)})
			dithered.SetGray(i+1, j+1, color.Gray{i16ToUI8(int16(dithered.GrayAt(i+1, j+1).Y) + quant)})
		}
	}
}
