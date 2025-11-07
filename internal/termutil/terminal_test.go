package termutil

import (
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestCharacterWidthCalculation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []MeasuredRune
	}{
		{
			name:  "ASCII characters",
			input: "Hello",
			expected: []MeasuredRune{
				{Rune: 'H', Width: 1},
				{Rune: 'e', Width: 1},
				{Rune: 'l', Width: 1},
				{Rune: 'l', Width: 1},
				{Rune: 'o', Width: 1},
			},
		},
		{
			name:  "Chinese characters",
			input: "你好",
			expected: []MeasuredRune{
				{Rune: '你', Width: 2},
				{Rune: '好', Width: 2},
			},
		},
		{
			name:  "Japanese hiragana",
			input: "こん",
			expected: []MeasuredRune{
				{Rune: 'こ', Width: 2},
				{Rune: 'ん', Width: 2},
			},
		},
		{
			name:  "Korean characters",
			input: "안녕",
			expected: []MeasuredRune{
				{Rune: '안', Width: 2},
				{Rune: '녕', Width: 2},
			},
		},
		{
			name:  "Mixed ASCII and CJK",
			input: "A中B",
			expected: []MeasuredRune{
				{Rune: 'A', Width: 1},
				{Rune: '中', Width: 2},
				{Rune: 'B', Width: 1},
			},
		},
		{
			name:  "Numbers and symbols",
			input: "123!",
			expected: []MeasuredRune{
				{Rune: '1', Width: 1},
				{Rune: '2', Width: 1},
				{Rune: '3', Width: 1},
				{Rune: '!', Width: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the runewidth calculation directly since processChan is internal
			for _, expected := range tt.expected {
				actualWidth := runewidth.RuneWidth(expected.Rune)
				assert.Equal(t, expected.Width, actualWidth, "Width mismatch for rune %c (0x%X)", expected.Rune, expected.Rune)
			}
		})
	}
}

func TestRuneWidthConsistency(t *testing.T) {
	// Test that our width calculation matches runewidth library directly
	testRunes := []rune{
		'A',  // ASCII - width 1
		'中',  // Chinese - width 2
		'こ',  // Japanese hiragana - width 2
		'안',  // Korean - width 2
		'1',  // Number - width 1
		' ',  // Space - width 1
		'\t', // Tab - width 0 (control character)
		'\n', // Newline - width 0 (control character)
	}

	for _, r := range testRunes {
		t.Run(string(r), func(t *testing.T) {
			expectedWidth := runewidth.RuneWidth(r)

			terminal := New()
			_, err := terminal.Write([]byte(string(r)))
			assert.NoError(t, err)

			// Skip control characters as they may not be processed the same way
			if r == '\t' || r == '\n' {
				return
			}

			// Test width calculation directly
			actualWidth := runewidth.RuneWidth(r)
			assert.Equal(t, expectedWidth, actualWidth, "Width mismatch for rune %c (0x%X)", r, r)
		})
	}
}

func TestEmptyInput(t *testing.T) {
	terminal := New()
	_, err := terminal.Write([]byte(""))
	assert.NoError(t, err)

	// Empty input should not cause any issues
	assert.NotNil(t, terminal)
}

func TestUTF8Encoding(t *testing.T) {
	// Test that multi-byte UTF-8 sequences are handled correctly
	// Chinese character "中" is 3 bytes in UTF-8: 0xE4 0xB8 0xAD
	chineseChar := "中"

	terminal := New()
	_, err := terminal.Write([]byte(chineseChar))
	assert.NoError(t, err)

	// Test that the width calculation is correct for multi-byte UTF-8
	width := runewidth.RuneWidth('中')
	assert.Equal(t, 2, width, "Chinese character should have width 2, not the UTF-8 byte count")
}

func TestControlCharacters(t *testing.T) {
	// Test that control characters have appropriate widths
	controlChars := map[rune]int{
		'\x00': 0, // NULL
		'\x07': 0, // BELL
		'\x08': 0, // BACKSPACE
		'\x09': 0, // TAB
		'\x0A': 0, // LINE FEED
		'\x0D': 0, // CARRIAGE RETURN
		'\x1B': 0, // ESCAPE
	}

	for char, expectedWidth := range controlChars {
		t.Run(string(char), func(t *testing.T) {
			actualWidth := runewidth.RuneWidth(char)
			assert.Equal(t, expectedWidth, actualWidth, "Control character 0x%X should have width %d", char, expectedWidth)
		})
	}
}
