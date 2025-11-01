package termutil

import (
	"image"
	"math"
	"strings"

	"github.com/zautumnz/tg/internal/app/tg/sixel"
)

type Sixel struct {
	X      uint16
	Y      uint64 // raw line
	Width  uint64
	Height uint64
	Image  image.Image
}

type VisibleSixel struct {
	ViewLineOffset int
	Sixel          Sixel
}

func (buffer *Buffer) addSixel(img image.Image, widthCells int, heightCells int) {
	buffer.sixels = append(buffer.sixels, Sixel{
		X:      buffer.CursorColumn(),
		Y:      buffer.cursorPosition.Line,
		Width:  uint64(widthCells),
		Height: uint64(heightCells),
		Image:  img,
	})
	if buffer.modes.SixelScrolling {
		buffer.cursorPosition.Line += uint64(heightCells)
	}
}

func (buffer *Buffer) clearSixelsAtRawLine(rawLine uint64) {
	var filtered []Sixel

	for _, sixelImage := range buffer.sixels {
		if sixelImage.Y+sixelImage.Height-1 >= rawLine && sixelImage.Y <= rawLine {
			continue
		}

		filtered = append(filtered, sixelImage)
	}

	buffer.sixels = filtered
}

func (buffer *Buffer) GetVisibleSixels() []VisibleSixel {

	firstLine := buffer.convertViewLineToRawLine(0)
	lastLine := buffer.convertViewLineToRawLine(buffer.viewHeight - 1)

	var visible []VisibleSixel

	for _, sixelImage := range buffer.sixels {
		if sixelImage.Y+sixelImage.Height-1 < firstLine {
			continue
		}
		if sixelImage.Y > lastLine {
			continue
		}

		visible = append(visible, VisibleSixel{
			ViewLineOffset: int(sixelImage.Y) - int(firstLine),
			Sixel:          sixelImage,
		})
	}

	return visible
}

func (t *Terminal) handleSixel(readChan chan MeasuredRune) (renderRequired bool) {

	var data []rune

	var inEscape bool

	for {
		r := <-readChan

		switch r.Rune {
		case 0x1b:
			inEscape = true
			continue
		case 0x5c:
			if inEscape {
				img, err := sixel.Decode(strings.NewReader(string(data)), t.theme.DefaultBackground())
				if err != nil {
					return false
				}
				w, h := t.windowManipulator.CellSizeInPixels()
				cw := int(math.Ceil(float64(img.Bounds().Dx()) / float64(w)))
				ch := int(math.Ceil(float64(img.Bounds().Dy()) / float64(h)))
				t.activeBuffer.addSixel(img, cw, ch)
				return true
			}
		}

		inEscape = false

		data = append(data, r.Rune)
	}
}
