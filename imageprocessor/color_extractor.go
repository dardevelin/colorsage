package imageprocessor

import (
	"image"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorExtractor processor extracts the color frequencies
type ColorExtractor struct{}

func (ce ColorExtractor) Name() string {
	return "ColorExtractor"
}

func (ce ColorExtractor) Process(img image.Image) (map[string]int, error) {
	colorMap := make(map[colorful.Color]int)

	// Iterate over each pixel in the image
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := colorful.Color{
				R: float64(r) / 65535.0,
				G: float64(g) / 65535.0,
				B: float64(b) / 65535.0,
			}
			colorMap[c]++
		}
	}

	// Convert colorful.Color map to hex string map
	hexMap := make(map[string]int)
	for c, count := range colorMap {
		hex := c.Hex()
		hexMap[hex] = count
	}

	return hexMap, nil
}
