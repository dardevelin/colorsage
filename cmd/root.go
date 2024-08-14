package cmd

import (
	"colorsage/imageprocessor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filePaths []string
var sequential bool
var quantizerType string
var fast bool
var all bool
var rawOutput bool
var includeFullColorExtract bool

var rootCmd = &cobra.Command{
	Use:   "colorsage [files...]",
	Short: "Process images and extract color palettes using various quantization algorithms.",
	Long: `colorsage is a tool for analyzing images and extracting their color palettes.
It supports multiple quantization algorithms to generate a reduced color palette.

By default, the tool runs the fastest quantizer (KMeansQuantizer).
You can use the --all flag to run all available quantizers (KMeans, MedianCut, and Average), 
or specify a particular quantizer using the --quantizer flag.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePaths = args // Capture the file paths from command-line arguments

		processors := []imageprocessor.ImageProcessor{
			&imageprocessor.ColorExtractor{},
		}

		var quantizers []imageprocessor.Quantizer

		// Determine which quantizers to run based on the command-line arguments
		if all {
			quantizers = append(quantizers, imageprocessor.KMeansQuantizer{})
			quantizers = append(quantizers, imageprocessor.MedianCutQuantizer{})
			quantizers = append(quantizers, imageprocessor.AverageQuantizer{})
		} else if quantizerType != "" {
			// Run only the specified quantizer
			quantizer, err := getQuantizerByName(quantizerType)
			if err != nil {
				fmt.Println(err)
				return
			}
			quantizers = append(quantizers, quantizer)
		} else {
			// Run only the fastest quantizer if neither --all nor --quantizer are specified
			quantizers = append(quantizers, imageprocessor.KMeansQuantizer{})
		}

		// Process the images using the selected quantizers
		results := imageprocessor.ProcessPipeline(filePaths, processors, quantizers, sequential)

		// Display results in a table format
		DisplayResultsTable(results)

		// Write all quantizer outputs to a file
		WriteResultsToFile("colors.txt", results)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&sequential, "sequential", "s", false, "Run the image processing pipeline sequentially (default: parallel)")
	rootCmd.PersistentFlags().StringVarP(&quantizerType, "quantizer", "q", "", "Specify which quantizer to use (kmeans, mediancut, average). If not specified, the fastest quantizer is used.")
	rootCmd.PersistentFlags().BoolVarP(&fast, "fast", "f", false, "Run only the fastest quantizer (default: KMeansQuantizer). This flag overrides running all quantizers.")
	rootCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "Run all available quantizers (KMeans, MedianCut, Average).")
	rootCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "Output raw results without UI elements, suitable for piping or redirection.")
	rootCmd.PersistentFlags().BoolVar(&includeFullColorExtract, "full-color-extract", false, "Include full color extraction details in the output.")
}

// getQuantizerByName returns the quantizer instance based on the provided name
func getQuantizerByName(name string) (imageprocessor.Quantizer, error) {
	switch name {
	case "kmeans":
		return imageprocessor.KMeansQuantizer{}, nil
	case "mediancut":
		return imageprocessor.MedianCutQuantizer{}, nil
	case "average":
		return imageprocessor.AverageQuantizer{}, nil
	default:
		return nil, fmt.Errorf("invalid quantizer type: %s. Supported types: kmeans, mediancut, average", name)
	}
}
