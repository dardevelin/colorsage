package imageprocessor

import (
	"math"
	"math/rand"

	"github.com/lucasb-eyer/go-colorful"
)

type KMeansQuantizer struct{}

func (q KMeansQuantizer) Name() string {
	return "KMeansQuantizer"
}

func (q KMeansQuantizer) Quantize(colorMap map[string]int, numColors int) map[string]int {
	colors := q.extractColors(colorMap)
	clusters := q.kmeans(colors, numColors)

	quantizedPalette := make(map[string]int)
	for _, cluster := range clusters {
		centroid := q.clusterCentroid(cluster)
		quantizedPalette[centroid.Hex()] = len(cluster)
	}

	return quantizedPalette
}

func (q KMeansQuantizer) extractColors(colorMap map[string]int) []colorful.Color {
	var colors []colorful.Color
	for hex := range colorMap {
		color, err := colorful.Hex(hex)
		if err == nil {
			colors = append(colors, color)
		}
	}
	return colors
}

func (q KMeansQuantizer) kmeans(colors []colorful.Color, numClusters int) [][]colorful.Color {
	centroids := make([]colorful.Color, numClusters)
	for i := range centroids {
		centroids[i] = colors[rand.Intn(len(colors))]
	}

	clusters := make([][]colorful.Color, numClusters)

	for i := 0; i < 10; i++ {
		for j := range clusters {
			clusters[j] = nil
		}

		for _, color := range colors {
			nearest := 0
			minDistance := math.MaxFloat64
			for k, centroid := range centroids {
				distance := color.DistanceLab(centroid)
				if distance < minDistance {
					nearest = k
					minDistance = distance
				}
			}
			clusters[nearest] = append(clusters[nearest], color)
		}

		for k := range centroids {
			centroids[k] = q.clusterCentroid(clusters[k])
		}
	}

	return clusters
}

func (q KMeansQuantizer) clusterCentroid(cluster []colorful.Color) colorful.Color {
	var r, g, b float64
	for _, color := range cluster {
		r += color.R
		g += color.G
		b += color.B
	}
	n := float64(len(cluster))
	return colorful.Color{
		R: r / n,
		G: g / n,
		B: b / n,
	}
}
