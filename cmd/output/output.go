// cmd/output/output.go
package output

import (
	"colorsage/config"
	"colorsage/imageprocessor"
)

// DisplayResults coordinates the display of results, choosing the appropriate output methods
func DisplayResults(results []imageprocessor.ImageResult) {
	if config.RawOutput || !IsOutputTerminal() {
		displayRawResults(results)
	} else {
		displayPrettyResults(results)
	}
}

// WriteResultsToFile coordinates the writing of results to a file, choosing the appropriate output methods
func WriteResultsToFile(filename string, results []imageprocessor.ImageResult) {
	if config.RawOutput || !IsOutputTerminal() {
		writeRawResultsToFile(filename, results)
	} else {
		writePrettyResultsToFile(filename, results)
	}
}
