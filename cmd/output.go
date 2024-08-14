package cmd

import (
	"colorsage/imageprocessor"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
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

// displayPrettyResults shows results in a nice table format with colors and generates palette images
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

		// Print the results for each quantizer and generate the palette image
		for _, quantizerName := range []string{"KMeansQuantizer", "MedianCutQuantizer", "AverageQuantizer"} {
			if palette, ok := result.Results[quantizerName]; ok {
				for hex, count := range palette {
					table.Append([]string{"", quantizerName, fmt.Sprintf(BackgroundColor(hex)+"%s"+Reset, hex), fmt.Sprintf("%d", count)})
				}

				// Generate an image for the palette
				imageFilename := generatePaletteFilename(result.FilePath, quantizerName)
				err := GeneratePaletteImage(palette, imageFilename)
				if err != nil {
					fmt.Printf(Red+"❌ Error generating palette image: %v"+Reset+"\n", err)
				} else {
					fmt.Printf(Cyan+"Image saved as: %s"+Reset+"\n", imageFilename)
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

// GeneratePaletteImage generates an image showing the color palette
func GeneratePaletteImage(colors map[string]int, filename string) error {
	const blockWidth, blockHeight = 50, 50
	imgWidth := len(colors) * blockWidth
	imgHeight := blockHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	i := 0
	for hex := range colors {
		c, err := colorful.Hex(hex)
		if err != nil {
			return err
		}
		draw.Draw(img, image.Rect(i*blockWidth, 0, (i+1)*blockWidth, blockHeight), &image.Uniform{C: color.RGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 255}}, image.Point{}, draw.Src)
		i++
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// generatePaletteFilename creates a filename for the palette image based on the file and quantizer
func generatePaletteFilename(filePath, quantizerName string) string {
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	return fmt.Sprintf("%s_%s_palette.png", nameWithoutExt, quantizerName)
}

// displayColorBlocks prints the color blocks in the terminal
func displayColorBlocks(colors map[string]int) {
	for hex := range colors {
		// Print a color block using ANSI background color codes
		fmt.Printf(BackgroundColor(hex) + "     " + Reset)
	}
	fmt.Println()
}
