package main

import (
	"bufio"
	"image"
	"image/draw"
	"os"
)

import (
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ImageFromFile(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	img, _, err := image.Decode(br)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func RGBAFromImage(img image.Image) (*image.RGBA, error) {
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m, nil
}
