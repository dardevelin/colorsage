package imageprocessor

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

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
