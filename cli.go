package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// ParseArgs parses command-line arguments and returns a slice of file paths and a flag for sequential processing
func ParseArgs() ([]string, bool) {
	sequential := flag.Bool("sequential", false, "Run the image processing pipeline sequentially. (default: parallel)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] <FILES>...\n", filepath.Base(os.Args[0]))
		fmt.Println("\nProcess images and extract color palettes.")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nArguments:")
		fmt.Println("  <FILES>  One or more image files to process. You can use glob patterns like '*.jpg'.")
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Collect all file paths passed as arguments
	files := flag.Args()

	// Ensure all files exist
	var validFiles []string
	for _, filePath := range files {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist.\n", filePath)
		} else {
			validFiles = append(validFiles, filePath)
		}
	}

	if len(validFiles) == 0 {
		fmt.Println("No valid files provided.")
		os.Exit(1)
	}

	return validFiles, *sequential
}
