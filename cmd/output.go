// cmd/output.go
package cmd

import (
	"colorsage/imageprocessor"
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

func BackgroundColor(hex string) string {
	color, err := colorful.Hex(hex)
	if err != nil {
		return ""
	}
	r, g, b := color.RGB255()
	return fmt.Sprintf("\033[48;2;%d;%d;%d;30m", r, g, b) // Use a lighter text color for contrast
}

func PrintResults(results []imageprocessor.ImageResult) {
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
				fmt.Printf("    - Color "+fmt.Sprintf("%s%s%s", BackgroundColor(hex), hex, Reset)+": %d occurrences\n", count)
			}
		}
		fmt.Println()
	}
}
