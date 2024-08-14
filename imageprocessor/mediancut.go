package imageprocessor

import (
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

type MedianCutQuantizer struct{}

func (q MedianCutQuantizer) Name() string {
	return "MedianCutQuantizer"
}

func (q MedianCutQuantizer) Quantize(colorMap map[string]int, numColors int) map[string]int {
	colors := q.extractColors(colorMap)
	if len(colors) == 0 {
		return map[string]int{}
	}

	initialBox := colorBox{colors: colors}
	boxes := q.medianCut([]colorBox{initialBox}, numColors)

	quantizedPalette := make(map[string]int)
	for _, box := range boxes {
		centroid := q.boxCentroid(box)
		quantizedPalette[centroid.Hex()] = len(box.colors)
	}

	return quantizedPalette
}

type colorBox struct {
	colors []colorful.Color
}

func (q MedianCutQuantizer) extractColors(colorMap map[string]int) []colorful.Color {
	var colors []colorful.Color
	for hex := range colorMap {
		color, err := colorful.Hex(hex)
		if err == nil {
			colors = append(colors, color)
		}
	}
	return colors
}

func (q MedianCutQuantizer) medianCut(boxes []colorBox, numColors int) []colorBox {
	for len(boxes) < numColors {
		newBoxes := []colorBox{}
		for _, box := range boxes {
			if len(box.colors) <= 1 {
				newBoxes = append(newBoxes, box)
				continue
			}

			rRange, gRange, bRange := q.colorRange(box.colors)
			var splitDim int
			if rRange >= gRange && rRange >= bRange {
				splitDim = 0
			} else if gRange >= rRange && gRange >= bRange {
				splitDim = 1
			} else {
				splitDim = 2
			}

			sort.Slice(box.colors, func(i, j int) bool {
				switch splitDim {
				case 0:
					return box.colors[i].R < box.colors[j].R
				case 1:
					return box.colors[i].G < box.colors[j].G
				case 2:
					return box.colors[i].B < box.colors[j].B
				default:
					return false
				}
			})

			median := len(box.colors) / 2
			box1 := colorBox{colors: box.colors[:median]}
			box2 := colorBox{colors: box.colors[median:]}

			newBoxes = append(newBoxes, box1, box2)
		}
		boxes = newBoxes
	}

	return boxes
}

func (q MedianCutQuantizer) colorRange(colors []colorful.Color) (float64, float64, float64) {
	var rMin, rMax, gMin, gMax, bMin, bMax float64
	rMin, gMin, bMin = 1.0, 1.0, 1.0

	for _, c := range colors {
		if c.R < rMin {
			rMin = c.R
		}
		if c.R > rMax {
			rMax = c.R
		}
		if c.G < gMin {
			gMin = c.G
		}
		if c.G > gMax {
			gMax = c.G
		}
		if c.B < bMin {
			bMin = c.B
		}
		if c.B > bMax {
			bMax = c.B
		}
	}

	return rMax - rMin, gMax - gMin, bMax - bMin
}

func (q MedianCutQuantizer) boxCentroid(box colorBox) colorful.Color {
	var r, g, b float64
	for _, color := range box.colors {
		r += color.R
		g += color.G
		b += color.B
	}
	n := float64(len(box.colors))
	return colorful.Color{
		R: r / n,
		G: g / n,
		B: b / n,
	}
}
