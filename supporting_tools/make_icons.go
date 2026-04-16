package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	createIcon("red.png", color.RGBA{255, 0, 0, 255})
	createIcon("yellow.png", color.RGBA{255, 255, 0, 255})
	createIcon("green.png", color.RGBA{0, 255, 0, 255})
	createIcon("blue.png", color.RGBA{0, 0, 255, 255})
	log.Println("Icons generated successfully!")
}

func createIcon(filename string, c color.Color) {
	// Create a 32x32 image
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Draw a solid block of color in the center, leaving a transparent border
	for x := 4; x < 28; x++ {
		for y := 4; y < 28; y++ {
			img.Set(x, y, c)
		}
	}

	// Save the file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", filename, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatalf("Failed to encode %s: %v", filename, err)
	}
}
