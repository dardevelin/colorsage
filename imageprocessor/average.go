package imageprocessor

import (
	"github.com/lucasb-eyer/go-colorful"
)

type AverageQuantizer struct{}

func (q AverageQuantizer) Name() string {
	return "AverageQuantizer"
}

func (q AverageQuantizer) Quantize(colorMap map[string]int, numColors int) map[string]int {
	colors := q.extractColors(colorMap)
	return q.simpleAverage(colors, numColors)
}

func (q AverageQuantizer) extractColors(colorMap map[string]int) []colorful.Color {
	var colors []colorful.Color
	for hex := range colorMap {
		color, err := colorful.Hex(hex)
		if err == nil {
			colors = append(colors, color)
		}
	}
	return colors
}

func (q AverageQuantizer) simpleAverage(colors []colorful.Color, numColors int) map[string]int {
	if len(colors) == 0 {
		return map[string]int{}
	}

	// Divide the colors into equal-sized buckets and average each bucket
	bucketSize := len(colors) / numColors
	if bucketSize == 0 {
		bucketSize = 1
	}

	quantizedPalette := make(map[string]int)

	for i := 0; i < len(colors); i += bucketSize {
		end := i + bucketSize
		if end > len(colors) {
			end = len(colors)
		}
		bucket := colors[i:end]
		centroid := q.bucketCentroid(bucket)
		quantizedPalette[centroid.Hex()] = len(bucket)
	}

	return quantizedPalette
}

func (q AverageQuantizer) bucketCentroid(bucket []colorful.Color) colorful.Color {
	var r, g, b float64
	for _, color := range bucket {
		r += color.R
		g += color.G
		b += color.B
	}
	n := float64(len(bucket))
	return colorful.Color{
		R: r / n,
		G: g / n,
		B: b / n,
	}
}
