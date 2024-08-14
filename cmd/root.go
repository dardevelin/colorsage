// cmd/root.go
package cmd

import (
	"colorsage/cmd/output"
	"colorsage/config"
	"colorsage/imageprocessor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
		config.FilePaths = args // Capture the file paths from command-line arguments

		processors := []imageprocessor.ImageProcessor{
			&imageprocessor.ColorExtractor{},
		}

		var quantizers []imageprocessor.Quantizer

		// Determine which quantizers to run based on the command-line arguments
		if config.All {
			quantizers = append(quantizers, imageprocessor.KMeansQuantizer{})
			quantizers = append(quantizers, imageprocessor.MedianCutQuantizer{})
			quantizers = append(quantizers, imageprocessor.AverageQuantizer{})
		} else if config.QuantizerType != "" {
			// Run only the specified quantizer
			quantizer, err := getQuantizerByName(config.QuantizerType)
			if err != nil {
				fmt.Println(err)
				return
			}
			quantizers = append(quantizers, quantizer)
		} else if config.Fast {
			// Run only the fastest quantizer if --fast is specified
			quantizers = append(quantizers, imageprocessor.KMeansQuantizer{})
		} else {
			// Run only the fastest quantizer if neither --all nor --quantizer nor --fast are specified
			quantizers = append(quantizers, imageprocessor.KMeansQuantizer{})
		}

		// Process the images using the selected quantizers
		results := imageprocessor.ProcessPipeline(config.FilePaths, processors, quantizers, config.Sequential)

		// Generate palette images if the flag is set
		if config.GeneratePaletteImagesInCurrentDir {
			for _, result := range results {
				for quantizerName, palette := range result.Results {
					filePath := output.GeneratePaletteFilename(result.FilePath, quantizerName)
					err := output.GeneratePaletteImage(palette, filePath)
					if err != nil {
						fmt.Printf("Error generating palette image for %s: %v\n", result.FilePath, err)
					}
				}
			}
		}

		// Display results in a table format
		output.DisplayResults(results)

		// Write all quantizer outputs to a file
		output.WriteResultsToFile("colors.txt", results)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Sequential, "sequential", "s", false, "Run the image processing pipeline sequentially (default: parallel)")
	rootCmd.PersistentFlags().StringVarP(&config.QuantizerType, "quantizer", "q", "", "Specify which quantizer to use (kmeans, mediancut, average). If not specified, the fastest quantizer is used.")
	rootCmd.PersistentFlags().BoolVar(&config.Fast, "fast", false, "Run only the fastest quantizer (default: KMeansQuantizer). This flag overrides running all quantizers.")
	rootCmd.PersistentFlags().BoolVar(&config.All, "all", false, "Run all available quantizers (KMeans, MedianCut, Average).")
	rootCmd.PersistentFlags().BoolVar(&config.RawOutput, "raw", false, "Output raw results without UI elements, suitable for piping or redirection.")
	rootCmd.PersistentFlags().BoolVar(&config.IncludeFullColorExtract, "full", false, "Include full color extraction details in the output.")
	rootCmd.PersistentFlags().BoolVar(&config.GeneratePaletteImagesInCurrentDir, "generate-palette-images-in-current-dir", false, "Generate palette images in the current directory instead of alongside the image files.")
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
