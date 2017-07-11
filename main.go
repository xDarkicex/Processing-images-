package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"path"
	"time"

	"image/jpeg"

	"runtime"

	"log"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/artyom/smartcrop"
	"github.com/fogleman/primitive/primitive"
	"github.com/pkg/errors"
)

func init() {

}

func main() {
	img, err := openImg("original.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img, err = crop(img, 1000, 1000)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImg(img, ".", "cropped.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img, err = openImg("cropped.jpg")
	if err != nil {
		log.Fatal(err)
	}
	sat := saturate(img)
	err = saveImg(sat, ".", "saturated.jpg")
	if err != nil {
		log.Fatal(err)
	}
	shrp := sharpen(sat)
	err = saveImg(shrp, ".", "sharpen.jpg")
	if err != nil {
		log.Fatal(err)
	}
	pri := primitivePicture(sat)
	err = saveImg(pri, ".", "primative.jpg")
	if err != nil {
		log.Fatal(err)
	}

}

func openImg(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot open "+path)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "Decoding the image failed.")
	}
	return img, nil
}

func saveImg(img image.Image, pname, fname string) error {
	fpath := path.Join(pname, fname)
	f, err := os.Create(fpath)
	if err != nil {
		return errors.Wrap(err, "Cannot create file: "+fpath)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return errors.Wrap(err, "Failed to encode the image as jpeg.")
	}
	return nil
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func crop(img image.Image, width, height int) (image.Image, error) {
	r, err := smartcrop.Crop(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "Smartcrop Failed!")
	}
	si, ok := (img).(SubImager)
	if !ok {
		return nil, errors.New("Crop(): img does not support SubImage()")
	}
	subImg := si.SubImage(r)
	return subImg, nil
}

func saturate(img image.Image) image.Image {
	return adjust.Saturation(img, 0.5)
}

func multiply(img image.Image) image.Image {
	return blend.Multiply(img, img)
}

func sharpen(img image.Image) image.Image {
	return effect.UnsharpMask(img, 0.6, 1.2)
}

func primitivePicture(img image.Image) image.Image {
	img = transform.Resize(img, 256, 256, transform.Linear)
	rand.Seed(time.Now().UTC().UnixNano())
	bg := primitive.MakeColor(primitive.AverageImageColor(img))
	model := primitive.NewModel(img, bg, 1024, runtime.NumCPU())
	for i := 0; i < 100; i++ {
		fmt.Print(".")
		model.Step(primitive.ShapeType(5), 128, 0)
	}
	return model.Context.Image()
}
