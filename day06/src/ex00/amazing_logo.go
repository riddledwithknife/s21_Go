package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	pink := color.RGBA{R: 255, G: 105, B: 180, A: 255}
	blue := color.RGBA{R: 30, G: 144, B: 255, A: 255}
	black := color.RGBA{A: 255}

	for x := 100; x < 200; x++ {
		for y := 100; y < 200; y++ {
			img.Set(x, y, pink)
		}
	}

	for x := 125; x <= 150; x++ {
		for y := 125; y <= 150; y++ {
			if (x == 125 && y == 125) || (x == 150 && y == 150) {
				img.Set(x, y, black)
			}
			if (x == 125 && y == 150) || (x == 150 && y == 125) {
				img.Set(x, y, black)
			}
		}
	}

	for x := 130; x <= 145; x++ {
		img.Set(x, 175, black)
	}

	for x := 90; x <= 210; x++ {
		for y := 70; y <= 100; y++ {
			if x >= 100 && x <= 200 && y >= 75 && y <= 95 {
				img.Set(x, y, blue)
			} else if (x == 90 || x == 210) && (y >= 80 && y <= 85) {
				img.Set(x, y, blue)
			}
		}
	}

	outputFile, err := os.Create("amazing_logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
