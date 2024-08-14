package main

import (
	"fmt"
	"image"
	"os"

	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format

	"github.com/lucasb-eyer/go-colorful"
)

// ImageProcessor interface for processing images
type ImageProcessor interface {
	Name() string
	Process(img image.Image) (map[string]int, error)
}

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

// ImageResult holds the results of processing an image
type ImageResult struct {
	FilePath string
	Results  map[string]map[string]int
	Err      error
}

// ProcessImage processes a single image through a pipeline of processors
func ProcessImage(filePath string, processors []ImageProcessor) ImageResult {
	file, err := os.Open(filePath)
	if err != nil {
		return ImageResult{FilePath: filePath, Err: err}
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return ImageResult{FilePath: filePath, Err: err}
	}

	results := make(map[string]map[string]int)
	for _, processor := range processors {
		result, err := processor.Process(img)
		if err != nil {
			return ImageResult{FilePath: filePath, Err: err}
		}
		results[processor.Name()] = result
	}

	return ImageResult{FilePath: filePath, Results: results}
}

// PrintResults prints the results of image processing
func PrintResults(results []ImageResult) {
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("❌ Error processing file %s: %v\n", result.FilePath, result.Err)
			continue
		}

		fmt.Printf("✅ Results for %s:\n", result.FilePath)
		for processorName, palette := range result.Results {
			fmt.Printf("  Processor: %s\n", processorName)
			for hex, count := range palette {
				fmt.Printf("    - Color %s: %d occurrences\n", hex, count)
			}
		}
		fmt.Println()
	}
}
