// cmd/output/text_output.go
package output

import (
	"colorsage/config"
	"colorsage/imageprocessor"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

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
		colorSummary := SummarizeColors(result.Results["ColorExtractor"])
		table.Append([]string{result.FilePath, "Summary", fmt.Sprintf("Total Colors: %d", colorSummary.TotalColors), ""})
		table.Append([]string{"", "Summary", fmt.Sprintf("Most Frequent: %s", colorSummary.MostFrequentColor), fmt.Sprintf("%d", colorSummary.MostFrequentCount)})
		table.Append([]string{"", "Summary", fmt.Sprintf("Least Frequent: %s", colorSummary.LeastFrequentColor), fmt.Sprintf("%d", colorSummary.LeastFrequentCount)})

		// Optionally, print full color extraction details
		if config.IncludeFullColorExtract {
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
		colorSummary := SummarizeColors(result.Results["ColorExtractor"])
		fmt.Printf("File: %s, Summary: Total Colors: %d\n", result.FilePath, colorSummary.TotalColors)
		fmt.Printf("Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount)
		fmt.Printf("Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount)

		// Optionally, print full color extraction details
		if config.IncludeFullColorExtract {
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

// writePrettyResultsToFile writes formatted results to a file
func writePrettyResultsToFile(filename string, results []imageprocessor.ImageResult) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	for _, result := range results {
		sb.WriteString(formatResultsForFile(result))
		sb.WriteString("\n")
	}

	file.WriteString(sb.String())
}

// writeRawResultsToFile writes raw results to a file
func writeRawResultsToFile(filename string, results []imageprocessor.ImageResult) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating results file:", err)
		return
	}
	defer file.Close()

	for _, result := range results {
		if result.Err != nil {
			fmt.Fprintf(file, "Error processing file %s: %v\n", result.FilePath, result.Err)
			continue
		}

		// Summary stats
		colorSummary := SummarizeColors(result.Results["ColorExtractor"])
		fmt.Fprintf(file, "File: %s, Summary: Total Colors: %d\n", result.FilePath, colorSummary.TotalColors)
		fmt.Fprintf(file, "Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount)
		fmt.Fprintf(file, "Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount)

		// Optionally, print full color extraction details
		if config.IncludeFullColorExtract {
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
	colorSummary := SummarizeColors(result.Results["ColorExtractor"])
	sb.WriteString(fmt.Sprintf("Results for %s:\n", result.FilePath))
	sb.WriteString(fmt.Sprintf("    - Total Colors: %d\n", colorSummary.TotalColors))
	sb.WriteString(fmt.Sprintf("    - Most Frequent: %s, Occurrences: %d\n", colorSummary.MostFrequentColor, colorSummary.MostFrequentCount))
	sb.WriteString(fmt.Sprintf("    - Least Frequent: %s, Occurrences: %d\n", colorSummary.LeastFrequentColor, colorSummary.LeastFrequentCount))

	// Optionally, print full color extraction details
	if config.IncludeFullColorExtract {
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
