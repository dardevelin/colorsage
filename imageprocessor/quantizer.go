package imageprocessor

import (
	"math"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

// Quantizer is responsible for reducing the number of colors in the image
type Quantizer struct{}

func (q Quantizer) Name() string {
	return "QuantizedPalette"
}

func (q Quantizer) Quantize(colorMap map[string]int, numColors int) map[string]int {
	// Step 1: Aggregate all colors into a slice
	var colors []colorful.Color
	var frequencies []int
	for hex, count := range colorMap {
		color, err := colorful.Hex(hex)
		if err != nil {
			continue
		}
		colors = append(colors, color)
		frequencies = append(frequencies, count)
	}

	// Step 2: Initialize the quantized color map
	quantizedColors := make(map[colorful.Color]int)

	// Step 3: For each color, find the nearest quantized color and update the frequency
	for i, color := range colors {
		nearestColor := findNearestColor(color, quantizedColors)
		quantizedColors[nearestColor] += frequencies[i]
	}

	// Step 4: Reduce to the most common colors if necessary
	if len(quantizedColors) > numColors {
		quantizedColors = reduceColors(quantizedColors, numColors)
	}

	// Step 5: Convert to hex map for output
	quantizedPalette := make(map[string]int)
	for color, count := range quantizedColors {
		quantizedPalette[color.Hex()] = count
	}

	return quantizedPalette
}

// findNearestColor finds the nearest color in the quantized color map
func findNearestColor(color colorful.Color, quantizedColors map[colorful.Color]int) colorful.Color {
	if len(quantizedColors) == 0 {
		return color
	}

	var nearestColor colorful.Color
	minDistance := math.MaxFloat64
	for quantizedColor := range quantizedColors {
		distance := color.DistanceLab(quantizedColor)
		if distance < minDistance {
			minDistance = distance
			nearestColor = quantizedColor
		}
	}

	return nearestColor
}

// reduceColors reduces the quantized color map to the most common colors
func reduceColors(quantizedColors map[colorful.Color]int, numColors int) map[colorful.Color]int {
	type kv struct {
		Key   colorful.Color
		Value int
	}
	var sortedColors []kv
	for k, v := range quantizedColors {
		sortedColors = append(sortedColors, kv{k, v})
	}
	sort.Slice(sortedColors, func(i, j int) bool {
		return sortedColors[i].Value > sortedColors[j].Value
	})

	// Keep only the most common colors
	reducedColors := make(map[colorful.Color]int)
	for i := 0; i < numColors && i < len(sortedColors); i++ {
		reducedColors[sortedColors[i].Key] = sortedColors[i].Value
	}

	return reducedColors
}
