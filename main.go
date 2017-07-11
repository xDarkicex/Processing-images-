package main

import (
	"image"
	"os"

	"image/jpeg"

	"github.com/pkg/errors"
)

func init() {

}

func main() {

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
