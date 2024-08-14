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

var rootCmd = &cobra.Command{
	Use:   "colorsage [files...]",
	Short: "Process images and extract color palettes using various quantization algorithms.",
	Long: `colorsage is a tool for analyzing images and extracting their color palettes.
It supports multiple quantization algorithms to generate a reduced color palette.

By default, the tool runs all available quantizers (KMeans, MedianCut, and Average) to give a comprehensive analysis.
You can use the --fast flag to run only the fastest quantizer (default: KMeansQuantizer), or specify a particular
quantizer using the --quantizer flag.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePaths = args // Capture the file paths from command-line arguments

		processors := []imageprocessor.ImageProcessor{
			&imageprocessor.ColorExtractor{},
		}

		// Select the quantizer based on the command-line argument
		var quantizer imageprocessor.Quantizer
		switch quantizerType {
		case "kmeans":
			quantizer = imageprocessor.KMeansQuantizer{}
		case "mediancut":
			quantizer = imageprocessor.MedianCutQuantizer{}
		case "average":
			quantizer = imageprocessor.AverageQuantizer{}
		case "":
			if fast {
				quantizer = imageprocessor.KMeansQuantizer{} // Assuming KMeans is the fastest
			} else {
				quantizer = imageprocessor.KMeansQuantizer{} // Default to KMeans if no other option is chosen
			}
		default:
			fmt.Println("Invalid quantizer type. Supported types: kmeans, mediancut, average")
			return
		}

		results := imageprocessor.ProcessPipeline(filePaths, processors, quantizer, sequential)

		PrintResults(results)
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
}
