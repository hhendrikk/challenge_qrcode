package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

const QRCODE_SIZE_PIXEL = 27

func main() {
	defer watch(time.Now(), "process time search qrcode")
	file, err := os.Open("./qrcodes.png")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	qrcodes, err := png.Decode(file)

	if err != nil {
		fmt.Println(err)
	}

	width := qrcodes.Bounds().Dx()
	newImage := image.NewNRGBA(image.Rect(0, 0, QRCODE_SIZE_PIXEL, QRCODE_SIZE_PIXEL))

	for x := 0; x < width; x += QRCODE_SIZE_PIXEL {
		for y := 0; y < width; y += QRCODE_SIZE_PIXEL {
			qrcode := crop(qrcodes, x, y, x+QRCODE_SIZE_PIXEL, y+QRCODE_SIZE_PIXEL)

			for i := 0; i < QRCODE_SIZE_PIXEL; i++ {
				for j := 0; j < QRCODE_SIZE_PIXEL; j++ {
					qrcodeRbga := color.NRGBAModel.Convert(qrcode.At(i, j)).(color.NRGBA)
					newImageRgba := color.NRGBAModel.Convert(newImage.At(i, j)).(color.NRGBA)

					r := qrcodeRbga.R ^ newImageRgba.R
					g := qrcodeRbga.G ^ newImageRgba.G
					b := qrcodeRbga.B ^ newImageRgba.B

					newImage.SetNRGBA(i, j, color.NRGBA{R: r, G: g, B: b, A: qrcodeRbga.A})
				}
			}
		}
	}

	fileResult, err := os.Create("./results.png")

	if err != nil {
		log.Fatalln(err)
	}

	err = png.Encode(fileResult, newImage)

	if err != nil {
		log.Fatalln(err)
	}
}

func crop(i image.Image, x int, y int, width int, height int) image.Image {
	cropSize := image.Rect(x, y, width, height)
	crop := i.(SubImager).SubImage(cropSize)
	newImage := image.NewNRGBA(image.Rect(0, 0, QRCODE_SIZE_PIXEL, QRCODE_SIZE_PIXEL))

	for i := x; i < width; i++ {
		for j := y; j < height; j++ {
			c := color.NRGBAModel.Convert(crop.At(i, j)).(color.NRGBA)
			newImage.SetNRGBA(i-x, j-y, c)
		}
	}

	return newImage
}

func watch(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
