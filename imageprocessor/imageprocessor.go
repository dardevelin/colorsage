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

	// Assuming the first processor is always a ColorExtractor
	if colorResults, ok := results["ColorExtractor"]; ok {
		quantizer := Quantizer{}
		quantizedPalette := quantizer.Quantize(colorResults, 5)
		results[quantizer.Name()] = quantizedPalette
	}

	return ImageResult{FilePath: filePath, Results: results}
}

// ProcessPipeline takes a list of file paths and processes them through the pipeline
func ProcessPipeline(filePaths []string, processors []ImageProcessor, sequential bool) []ImageResult {
	results := make([]ImageResult, len(filePaths))

	if sequential {
		fmt.Println(Yellow + "Running in sequential mode..." + Reset)
		for i, filePath := range filePaths {
			results[i] = ProcessImage(filePath, processors)
		}
	} else {
		fmt.Println(Yellow + "Running in parallel mode..." + Reset)
		var wg sync.WaitGroup
		wg.Add(len(filePaths))

		for i, filePath := range filePaths {
			go func(i int, filePath string) {
				defer wg.Done()
				results[i] = ProcessImage(filePath, processors)
			}(i, filePath)
		}

		wg.Wait()
	}

	return results
}

// PrintResults prints the results of image processing with colored output
func PrintResults(results []ImageResult) {
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf(Red+"❌ Error processing file %s: %v"+Reset+"\n", result.FilePath, result.Err)
			continue
		}

		fmt.Printf(Green+"✅ Results for %s:"+Reset+"\n", result.FilePath)
		for processorName, palette := range result.Results {
			if processorName == "QuantizedPalette" {
				fmt.Printf(Purple + "  Quantized Palette:" + Reset + "\n")
			} else {
				fmt.Printf(Cyan+"  Processor: %s"+Reset+"\n", processorName)
			}
			for hex, count := range palette {
				// Display the hex code with its corresponding color as a background
				fmt.Printf("    - Color "+fmt.Sprintf("%s%s%s", BackgroundColor(hex), hex, Reset)+": %d occurrences\n", count)
			}
		}
		fmt.Println()
	}
}