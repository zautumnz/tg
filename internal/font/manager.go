package font

import (
	_ "embed"
	"fmt"
	"image"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Style uint8

const (
	Regular Style = iota
)

type StyleName string

const (
	StyleRegular StyleName = "Regular"
)

type Manager struct {
	family       string
	regularFace  font.Face
	size         float64
	dpi          float64
	charSize     image.Point
	fontDotDepth int
}

func NewManager() *Manager {
	return &Manager{
		size: 16,
		dpi:  72,
	}
}

func (m *Manager) CharSize() image.Point {
	return m.charSize
}

func (m *Manager) IncreaseSize() {
	m.SetSize(m.size + 1)
}

func (m *Manager) DecreaseSize() {
	if m.size < 2 {
		return
	}
	m.SetSize(m.size - 1)
}

func (m *Manager) DotDepth() int {
	return m.fontDotDepth
}

func (m *Manager) DPI() float64 {
	return m.dpi
}

func (m *Manager) SetDPI(dpi float64) error {
	if dpi <= 0 {
		return fmt.Errorf("DPI must be >0")
	}
	m.dpi = dpi
	return nil
}

func (m *Manager) SetSize(size float64) error {
	m.size = size
	if m.regularFace != nil {
		// effectively reload fonts at new size
		m.SetFontByFamilyName(m.family)
	}
	return nil
}

func (m *Manager) createFace(f *opentype.Font) (font.Face, error) {
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    m.size,
		DPI:     m.dpi,
		Hinting: font.HintingFull,
	})
}

func (m *Manager) SetFontByFamilyName(name string) error {
	return m.loadDefaultFonts()
}

func (m *Manager) calcMetrics() error {

	face := m.regularFace

	var prevAdvance int
	for ch := rune(32); ch <= 126; ch++ {
		adv26, ok := face.GlyphAdvance(ch)
		if ok && adv26 > 0 {
			advance := int(adv26)
			if prevAdvance > 0 && prevAdvance != advance {
				return fmt.Errorf("the specified font is not monospaced: %d 0x%X=%d", prevAdvance, ch, advance)
			}
			prevAdvance = advance
		}
	}

	if prevAdvance == 0 {
		return fmt.Errorf("failed to calculate advance width for font face")
	}

	metrics := face.Metrics()

	m.charSize.X = int(math.Round(float64(prevAdvance) / m.dpi))
	m.charSize.Y = int(math.Round(float64(metrics.Height) / m.dpi))
	m.fontDotDepth = int(math.Round(float64(metrics.Ascent) / m.dpi))

	return nil
}

//go:embed SarasaTermCL-Regular.ttf
var sarasa []byte

func (m *Manager) loadDefaultFonts() error {

	regular, err := opentype.Parse(sarasa)
	if err != nil {
		return err
	}
	m.regularFace, err = m.createFace(regular)
	if err != nil {
		return err
	}

	return m.calcMetrics()
}

func (m *Manager) RegularFontFace() font.Face {
	return m.regularFace
}
