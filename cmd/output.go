package cmd

import (
	"colorsage/imageprocessor"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// DisplayResultsTable displays the results in a table format or raw, depending on the context
func DisplayResultsTable(results []imageprocessor.ImageResult) {
	if rawOutput || !IsOutputTerminal() {
		displayRawResults(results)
	} else {
		displayPrettyResults(results)
	}
}

// displayPrettyResults shows results in a nice table format with colors
func displayPrettyResults(results []imageprocessor.ImageResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Quantizer", "Color", "Occurrences"})

	for _, result := range results {
		if result.Err != nil {
			fmt.Printf(Red+"❌ Error processing file %s: %v"+Reset+"\n", result.FilePath, result.Err)
			continue
		}

		// Summary stats
		colorSummary := summarizeColors(result.Results["ColorExtractor"])
		table.Append([]string{result.FilePath, "Summary", fmt.Sprintf("Total Colors: %d", colorSummary.TotalColors), ""})
		table.Append([]string{"", "Summary", fmt.Sprintf("Most Frequent: %s", colorSummary.MostFrequentColor), fmt.Sprintf("%d", colorSummary.MostFrequentCount)})
		table.Append([]string{"", "Summary", fmt.Sprintf("Least Frequent: %s", colorSummary.LeastFrequentColor), fmt.Sprintf("%d", colorSummary.LeastFrequentCount)})

		// Optionally, print full color extraction details
		if includeFullColorExtract {
			if colorResults, ok := result.Results["ColorExtractor"]; ok {
				for hex, count := range colorResults {
					table.Append([]string{"", "ColorExtractor", fmt.Sprintf(BackgroundColor(hex)+"%s"+Reset, hex), fmt.Sprintf("%d", count)})
				}
			}
		}

		// Print the results for each quantizer
		for _, quantizerName := range []string{"KMeansQuantizer", "MedianCutQuantizer", "AverageQuantizer"} {
			if palette, ok := result.Results[quantizerName]; ok {
				for hex, count := range palette {
					table.Append([]string{"", quantizerName, fmt.Sprintf(BackgroundColor(hex)+"%s"+Reset, hex), fmt.Sprintf("%d", count)})
				}
			}
		}
	}
	table.SetRowLine(true)
	table.Render()
}

// displayRawResults outputs results in a simple format suitable for piping
func displayRawResults(results []imageprocessor.ImageResult) {
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("Error processing file %s: %v\n", result.FilePath, result.Err)
			continue
		}

		// Summary stats
		colorSummary := summarizeColors(result.Results["ColorExtractor"])
		fmt.Printf("File: %s, Summary: Total Colors: %d\n", result.FilePath, colorSummary.TotalColors)
		fmt.Printf("Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount)
		fmt.Printf("Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount)

		// Optionally, print full color extraction details
		if includeFullColorExtract {
			if colorResults, ok := result.Results["ColorExtractor"]; ok {
				for hex, count := range colorResults {
					fmt.Printf("File: %s, Quantizer: ColorExtractor, Color: %s, Occurrences: %d\n", result.FilePath, hex, count)
				}
			}
		}

		// Print the results for each quantizer
		for _, quantizerName := range []string{"KMeansQuantizer", "MedianCutQuantizer", "AverageQuantizer"} {
			if palette, ok := result.Results[quantizerName]; ok {
				for hex, count := range palette {
					fmt.Printf("File: %s, Quantizer: %s, Color: %s, Occurrences: %d\n", result.FilePath, quantizerName, hex, count)
				}
			}
		}
	}
}

// WriteResultsToFile writes all quantizer results to the specified file
func WriteResultsToFile(filename string, results []imageprocessor.ImageResult) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}
	defer file.Close()

	if rawOutput || !IsOutputTerminal() {
		writeRawResults(file, results)
	} else {
		writePrettyResults(file, results)
	}

	fmt.Println("All quantizer results have been written to", filename)
}

// writePrettyResults writes formatted results to a file
func writePrettyResults(file *os.File, results []imageprocessor.ImageResult) {
	var sb strings.Builder
	for _, result := range results {
		sb.WriteString(formatResultsForFile(result))
		sb.WriteString("\n")
	}

	file.WriteString(sb.String())
}

// writeRawResults writes raw results to a file
func writeRawResults(file *os.File, results []imageprocessor.ImageResult) {
	for _, result := range results {
		if result.Err != nil {
			fmt.Fprintf(file, "Error processing file %s: %v\n", result.FilePath, result.Err)
			continue
		}

		// Summary stats
		colorSummary := summarizeColors(result.Results["ColorExtractor"])
		fmt.Fprintf(file, "File: %s, Summary: Total Colors: %d\n", result.FilePath, colorSummary.TotalColors)
		fmt.Fprintf(file, "Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount)
		fmt.Fprintf(file, "Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount)

		// Optionally, print full color extraction details
		if includeFullColorExtract {
			if colorResults, ok := result.Results["ColorExtractor"]; ok {
				for hex, count := range colorResults {
					fmt.Fprintf(file, "File: %s, Quantizer: ColorExtractor, Color: %s, Occurrences: %d\n", result.FilePath, hex, count)
				}
			}
		}

		// Print the results for each quantizer
		for _, quantizerName := range []string{"KMeansQuantizer", "MedianCutQuantizer", "AverageQuantizer"} {
			if palette, ok := result.Results[quantizerName]; ok {
				for hex, count := range palette {
					fmt.Fprintf(file, "File: %s, Quantizer: %s, Color: %s, Occurrences: %d\n", result.FilePath, quantizerName, hex, count)
				}
			}
		}
	}
}

// formatResultsForFile formats the results for pretty file output
func formatResultsForFile(result imageprocessor.ImageResult) string {
	var sb strings.Builder
	if result.Err != nil {
		sb.WriteString(fmt.Sprintf("Error processing file %s: %v\n", result.FilePath, result.Err))
		return sb.String()
	}

	// Summary stats
	colorSummary := summarizeColors(result.Results["ColorExtractor"])
	sb.WriteString(fmt.Sprintf("Results for %s:\n", result.FilePath))
	sb.WriteString(fmt.Sprintf("    - Total Colors: %d\n", colorSummary.TotalColors))
	sb.WriteString(fmt.Sprintf("    - Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount))
	sb.WriteString(fmt.Sprintf("    - Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount))

	// Optionally, print full color extraction details
	if includeFullColorExtract {
		if colorResults, ok := result.Results["ColorExtractor"]; ok {
			for hex, count := range colorResults {
				sb.WriteString(fmt.Sprintf("    - Color %s: %d occurrences\n", hex, count))
			}
		}
	}

	// Then print the results for each quantizer
	for _, quantizerName := range []string{"KMeansQuantizer", "MedianCutQuantizer", "AverageQuantizer"} {
		if palette, ok := result.Results[quantizerName]; ok {
			sb.WriteString(fmt.Sprintf("Results for Quantizer: %s\n", quantizerName))
			for hex, count := range palette {
				sb.WriteString(fmt.Sprintf("    - Color %s: %d occurrences\n", hex, count))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// summarizeColors calculates summary statistics for a color palette
func summarizeColors(colorMap map[string]int) ColorSummary {
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
