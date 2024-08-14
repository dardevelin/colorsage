// cmd/output/terminal_output.go
package output

import (
	"fmt"
)

// DisplayColorBlocks prints the color blocks in the terminal
func DisplayColorBlocks(colors map[string]int) {
	for hex := range colors {
		// Print a color block using ANSI background color codes
		fmt.Printf(BackgroundColor(hex) + "     " + Reset)
	}
	fmt.Println()
}
