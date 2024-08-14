package imageprocessor

type Quantizer interface {
	Name() string
	Quantize(colorMap map[string]int, numColors int) map[string]int
}
