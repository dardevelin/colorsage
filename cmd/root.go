package cmd

import (
	"colorsage/imageprocessor"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var filePaths []string
var sequential bool
var quantizerType string
var fast bool
var all bool

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

		var allResults []string

		// Process the images using the selected quantizers
		for _, quantizer := range quantizers {
			results := imageprocessor.ProcessPipeline(filePaths, processors, quantizer, sequential)
			PrintResults(results)
			allResults = append(allResults, formatResultsForFile(results, quantizer.Name()))
		}

		// Write all quantizer outputs to a file with color
		writeResultsToFile("colors.txt", allResults)
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

// formatResultsForFile formats the results of processing for writing to a file with color
func formatResultsForFile(results []imageprocessor.ImageResult, quantizerName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(Purple+"Results for Quantizer: %s"+Reset+"\n", quantizerName))
	for _, result := range results {
		if result.Err != nil {
			sb.WriteString(fmt.Sprintf(Red+"❌ Error processing file %s: %v"+Reset+"\n", result.FilePath, result.Err))
			continue
		}
		sb.WriteString(fmt.Sprintf(Green+"✅ Results for %s:"+Reset+"\n", result.FilePath))
		for processorName, palette := range result.Results {
			if processorName == quantizerName {
				sb.WriteString(fmt.Sprintf(Cyan+"  Quantized Palette (%s):"+Reset+"\n", quantizerName))
				for hex, count := range palette {
					sb.WriteString(fmt.Sprintf("    - Color "+BackgroundColor(hex)+"%s"+Reset+": %d occurrences\n", hex, count))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// writeResultsToFile writes all quantizer results to the specified file with color
func writeResultsToFile(filename string, results []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(Red+"❌ Error creating results file:"+Reset, err)
		return
	}
	defer file.Close()

	for _, result := range results {
		file.WriteString(result)
		file.WriteString("\n")
	}

	fmt.Println(Green+"✅ All quantizer results have been written to", filename+Reset)
}
