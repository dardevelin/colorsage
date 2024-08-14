package imageprocessor

import (
	"image"
	"runtime"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorExtractor processor extracts the color frequencies
type ColorExtractor struct{}

func (ce ColorExtractor) Name() string {
	return "ColorExtractor"
}

func (ce ColorExtractor) Process(img image.Image) (map[string]int, error) {
	colorMap := make(map[colorful.Color]int)
	var mutex sync.Mutex

	// Determine the number of threads to use
	numThreads := runtime.NumCPU()
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	chunkHeight := height / numThreads

	var wg sync.WaitGroup

	// Process each chunk in parallel
	for i := 0; i < numThreads; i++ {
		wg.Add(1)

		go func(startY int) {
			defer wg.Done()

			localColorMap := make(map[colorful.Color]int)
			endY := startY + chunkHeight
			if i == numThreads-1 { // Ensure the last chunk covers the remainder of the image
				endY = height
			}

			// Iterate over the chunk's pixels
			for y := startY; y < endY; y++ {
				for x := 0; x < width; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					c := colorful.Color{
						R: float64(r) / 65535.0,
						G: float64(g) / 65535.0,
						B: float64(b) / 65535.0,
					}
					localColorMap[c]++
				}
			}

			// Safely merge local color map into the global color map
			mutex.Lock()
			for color, count := range localColorMap {
				colorMap[color] += count
			}
			mutex.Unlock()
		}(i * chunkHeight)
	}

	wg.Wait()

	// Convert colorful.Color map to hex string map
	hexMap := make(map[string]int)
	for c, count := range colorMap {
		hex := c.Hex()
		hexMap[hex] = count
	}

	return hexMap, nil
}
