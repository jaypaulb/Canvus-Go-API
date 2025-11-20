package canvus

import (
	"fmt"
	"regexp"
	"strings"
)

// Color helper functions for the Canvus API RRGGBBAA format.
//
// The Canvus API uses 8-character uppercase hex strings for colors:
// RRGGBBAA where RR=red, GG=green, BB=blue, AA=alpha (00=transparent, FF=opaque)

var colorRGBAPattern = regexp.MustCompile(`^[0-9A-F]{8}$`)

// ValidateColor checks if a string is a valid Canvus color (RRGGBBAA uppercase).
// Returns an error if the color format is invalid.
func ValidateColor(color string) error {
	if len(color) != 8 {
		return fmt.Errorf("color must be exactly 8 characters (RRGGBBAA), got %d", len(color))
	}
	if !colorRGBAPattern.MatchString(color) {
		return fmt.Errorf("color must be uppercase hex RRGGBBAA format, got %q", color)
	}
	return nil
}

// NormalizeColor converts a color string to the Canvus format (uppercase RRGGBBAA).
// Accepts:
//   - RRGGBBAA (already valid, just uppercased)
//   - RRGGBB (adds FF for opaque alpha)
//   - #RRGGBBAA or #RRGGBB (removes # prefix)
//
// Returns the normalized color or an error if the input format is invalid.
func NormalizeColor(color string) (string, error) {
	// Remove # prefix if present
	color = strings.TrimPrefix(color, "#")

	// Convert to uppercase
	color = strings.ToUpper(color)

	// Handle 6-character RGB (add opaque alpha)
	if len(color) == 6 {
		if matched, _ := regexp.MatchString(`^[0-9A-F]{6}$`, color); matched {
			return color + "FF", nil
		}
		return "", fmt.Errorf("invalid 6-character color format: %q", color)
	}

	// Validate 8-character RGBA
	if err := ValidateColor(color); err != nil {
		return "", err
	}

	return color, nil
}

// ColorToRGBA converts a Canvus color (RRGGBBAA) to separate R, G, B, A byte values.
// Returns an error if the color format is invalid.
func ColorToRGBA(color string) (r, g, b, a byte, err error) {
	if err := ValidateColor(color); err != nil {
		return 0, 0, 0, 0, err
	}

	var rr, gg, bb, aa uint
	_, err = fmt.Sscanf(color, "%02X%02X%02X%02X", &rr, &gg, &bb, &aa)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse color %q: %w", color, err)
	}

	return byte(rr), byte(gg), byte(bb), byte(aa), nil
}

// RGBAToColor converts separate R, G, B, A byte values to a Canvus color (RRGGBBAA).
func RGBAToColor(r, g, b, a byte) string {
	return fmt.Sprintf("%02X%02X%02X%02X", r, g, b, a)
}

// ColorToRGB converts a Canvus color (RRGGBBAA) to standard #RRGGBB format (discards alpha).
func ColorToRGB(color string) (string, error) {
	if err := ValidateColor(color); err != nil {
		return "", err
	}
	return "#" + color[:6], nil
}

// ColorWithAlpha creates a new color by setting the alpha channel of an existing color.
// Alpha should be 0-255 (0x00-0xFF).
func ColorWithAlpha(color string, alpha byte) (string, error) {
	if err := ValidateColor(color); err != nil {
		return "", err
	}
	return color[:6] + fmt.Sprintf("%02X", alpha), nil
}

// Common opaque colors for convenience
const (
	ColorBlack       = "000000FF"
	ColorWhite       = "FFFFFFFF"
	ColorRed         = "FF0000FF"
	ColorGreen       = "00FF00FF"
	ColorBlue        = "0000FFFF"
	ColorYellow      = "FFFF00FF"
	ColorCyan        = "00FFFFFF"
	ColorMagenta     = "FF00FFFF"
	ColorGray        = "808080FF"
	ColorLightGray   = "D3D3D3FF"
	ColorDarkGray    = "404040FF"
	ColorTransparent = "00000000"
)
