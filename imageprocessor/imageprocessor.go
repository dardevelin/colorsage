// imageprocessor/imageprocessor.go
package imageprocessor

import (
	"fmt"
	"image"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"
	"sync"

	_ "golang.org/x/image/bmp"  // Register BMP format
	_ "golang.org/x/image/tiff" // Register TIFF format
	_ "golang.org/x/image/webp" // Register WebP format
)

// ImageProcessor interface for processing images
type ImageProcessor interface {
	Name() string
	Process(img image.Image) (map[string]int, error)
}

// ImageResult holds the results of processing an image
type ImageResult struct {
	FilePath string
	Results  map[string]map[string]int
	Err      error
}

// ProcessImage processes a single image through a pipeline of processors and quantizers
func ProcessImage(filePath string, processors []ImageProcessor, quantizers []Quantizer) ImageResult {
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
	var colorPalette map[string]int

	// Step 1: Run ColorExtractor once
	for _, processor := range processors {
		if processor.Name() == "ColorExtractor" {
			colorPalette, err = processor.Process(img)
			if err != nil {
				return ImageResult{FilePath: filePath, Err: err}
			}
			results[processor.Name()] = colorPalette
			break // We only need to run ColorExtractor once
		}
	}

	// Step 2: Pass the extracted color palette to each quantizer
	numColors := 5 // Hardcoded number of colors, or make it configurable
	for _, quantizer := range quantizers {
		quantizedPalette, err := quantizer.Quantize(colorPalette, numColors)
		if err != nil {
			return ImageResult{FilePath: filePath, Err: err}
		}
		results[quantizer.Name()] = quantizedPalette
	}

	return ImageResult{FilePath: filePath, Results: results}
}

// ProcessPipeline takes a list of file paths and processes them through the pipeline
func ProcessPipeline(filePaths []string, processors []ImageProcessor, quantizers []Quantizer, sequential bool) []ImageResult {
	results := make([]ImageResult, len(filePaths))

	if sequential {
		fmt.Println(Yellow + "Running in sequential mode..." + Reset)
		for i, filePath := range filePaths {
			results[i] = ProcessImage(filePath, processors, quantizers)
		}
	} else {
		fmt.Println(Yellow + "Running in parallel mode..." + Reset)
		var wg sync.WaitGroup
		wg.Add(len(filePaths))

		for i, filePath := range filePaths {
			go func(i int, filePath string) {
				defer wg.Done()
				results[i] = ProcessImage(filePath, processors, quantizers)
			}(i, filePath)
		}

		wg.Wait()
	}

	return results
}
