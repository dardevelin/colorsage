package main

import (
	"colorsage/cli" // Import the package that handles command-line arguments
	"fmt"
	"sync"
)

// ProcessPipeline takes a list of file paths and processes them through the pipeline
func ProcessPipeline(filePaths []string, processors []ImageProcessor, sequential bool) []ImageResult {
	results := make([]ImageResult, len(filePaths))

	if sequential {
		fmt.Println("Running in sequential mode...")
		for i, filePath := range filePaths {
			results[i] = ProcessImage(filePath, processors)
		}
	} else {
		fmt.Println("Running in parallel mode...")
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

func main() {
	files, sequential := cli.ParseArgs() // Call the function to parse command-line arguments

	fmt.Println("Processing files:", files)
	fmt.Println()

	// Define the pipeline of processors
	processors := []ImageProcessor{
		ColorExtractor{},
		// Add more processors here if needed
	}

	// Process the files through the image processing pipeline
	results := ProcessPipeline(files, processors, sequential)

	// Print the results
	PrintResults(results)
}
