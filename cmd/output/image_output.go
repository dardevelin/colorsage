// cmd/output/image_output.go
package output

import (
	"colorsage/config"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// GeneratePaletteImage creates a PNG image representing the given colors and saves it to the specified file path.
func GeneratePaletteImage(colors map[string]int, filePath string) error {
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

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// GeneratePaletteFilename constructs the filename for the palette image based on the provided file path and quantizer name.
func GeneratePaletteFilename(filePath, quantizerName string) string {
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	if config.GeneratePaletteImagesInCurrentDir {
		return fmt.Sprintf("./%s_%s_palette.png", nameWithoutExt, quantizerName)
	}
	return fmt.Sprintf("%s_%s_palette.png", nameWithoutExt, quantizerName)
}
