// cmd/output/helpers.go
package output

import (
	"fmt"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)

// IsOutputTerminal checks if the output is a terminal
func IsOutputTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ANSI color codes for colorful output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

// BackgroundColor returns a string with the ANSI escape code to set the background color
func BackgroundColor(hex string) string {
	color, err := colorful.Hex(hex)
	if err != nil {
		return ""
	}
	r, g, b := color.RGB255()
	return fmt.Sprintf("\033[48;2;%d;%d;%d;30m", r, g, b) // Use a lighter text color for contrast
}

// SummarizeColors calculates summary statistics for a color palette
func SummarizeColors(colorMap map[string]int) ColorSummary {
	summary := ColorSummary{
		TotalColors:        len(colorMap),
		MostFrequentCount:  0,
		LeastFrequentCount: int(^uint(0) >> 1), // Initialize with max int
	}

	for hex, count := range colorMap {
		if count > summary.MostFrequentCount {
			summary.MostFrequentColor = hex
			summary.MostFrequentCount = count
		}
		if count < summary.LeastFrequentCount {
			summary.LeastFrequentColor = hex
			summary.LeastFrequentCount = count
		}
	}

	return summary
}

// ColorSummary holds summary statistics for a color palette
type ColorSummary struct {
	TotalColors        int
	MostFrequentColor  string
	MostFrequentCount  int
	LeastFrequentColor string
	LeastFrequentCount int
}
