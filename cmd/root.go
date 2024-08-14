package cmd

import (
	"colorsage/imageprocessor"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filePaths []string
var sequential bool

var rootCmd = &cobra.Command{
	Use:   "colorsage [files...]",
	Short: "Process images and extract color palettes",
	Long:  `colorsage is a tool for analyzing images and extracting their color palettes. It supports processing in parallel or sequentially.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePaths = args // Capture the file paths from command-line arguments

		processors := []imageprocessor.ImageProcessor{
			&imageprocessor.ColorExtractor{},
		}

		results := imageprocessor.ProcessPipeline(filePaths, processors, sequential)
		PrintResults(results)
	},
} // <- Closing brace instead of closing parenthesis

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&sequential, "sequential", "s", false, "Run the image processing pipeline sequentially (default: parallel)")
}
