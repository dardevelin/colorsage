package imageprocessor

// Quantizer interface for quantizing color palettes
type Quantizer interface {
	Name() string
	Quantize(colors map[string]int, numColors int) (map[string]int, error)
}
