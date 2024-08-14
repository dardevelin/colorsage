package cmd

import (
	"colorsage/imageprocessor" // Correct import path
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sequential bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "colorsage [files...]",
	Short: "Process images and extract color palettes",
	Long: `colorsage is a tool for analyzing images and extracting their color palettes.
It supports parallel and sequential processing of images.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Call the processing pipeline
		processors := []imageprocessor.ImageProcessor{
			imageprocessor.ColorExtractor{},
			// Add more processors here if needed
		}

		results := imageprocessor.ProcessPipeline(args, processors, sequential)
		imageprocessor.PrintResults(results)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.PersistentFlags().BoolVarP(&sequential, "sequential", "s", false, "Run the image processing pipeline sequentially (default: parallel)")
}
