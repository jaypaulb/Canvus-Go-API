package canvus

import "testing"

func TestValidateColor(t *testing.T) {
	tests := []struct {
		name    string
		color   string
		wantErr bool
	}{
		{"valid uppercase", "FF0000FF", false},
		{"valid black", "000000FF", false},
		{"valid transparent", "FFFFFF00", false},
		{"lowercase invalid", "ff0000ff", true},
		{"too short", "FF0000", true},
		{"too long", "FF0000FFFF", true},
		{"invalid chars", "GGGGGGGG", true},
		{"with hash", "#FF0000FF", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColor(tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColor(%q) error = %v, wantErr %v", tt.color, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"already valid", "FF0000FF", "FF0000FF", false},
		{"lowercase to uppercase", "ff0000ff", "FF0000FF", false},
		{"6-char RGB", "FF0000", "FF0000FF", false},
		{"6-char RGB lowercase", "ff0000", "FF0000FF", false},
		{"with hash 8-char", "#FF0000FF", "FF0000FF", false},
		{"with hash 6-char", "#FF0000", "FF0000FF", false},
		{"mixed case", "Ff00Ff80", "FF00FF80", false},
		{"invalid chars", "GGGGGG", "", true},
		{"too short", "FF00", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeColor(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorToRGBA(t *testing.T) {
	tests := []struct {
		name  string
		color string
		wantR byte
		wantG byte
		wantB byte
		wantA byte
	}{
		{"red", "FF0000FF", 255, 0, 0, 255},
		{"green", "00FF00FF", 0, 255, 0, 255},
		{"blue", "0000FFFF", 0, 0, 255, 255},
		{"white", "FFFFFFFF", 255, 255, 255, 255},
		{"black", "000000FF", 0, 0, 0, 255},
		{"transparent white", "FFFFFF00", 255, 255, 255, 0},
		{"50% gray", "80808080", 128, 128, 128, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, a, err := ColorToRGBA(tt.color)
			if err != nil {
				t.Fatalf("ColorToRGBA(%q) error = %v", tt.color, err)
			}
			if r != tt.wantR || g != tt.wantG || b != tt.wantB || a != tt.wantA {
				t.Errorf("ColorToRGBA(%q) = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					tt.color, r, g, b, a, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
			}
		})
	}
}

func TestRGBAToColor(t *testing.T) {
	tests := []struct {
		name string
		r, g, b, a byte
		want string
	}{
		{"red", 255, 0, 0, 255, "FF0000FF"},
		{"green", 0, 255, 0, 255, "00FF00FF"},
		{"blue", 0, 0, 255, 255, "0000FFFF"},
		{"white", 255, 255, 255, 255, "FFFFFFFF"},
		{"black", 0, 0, 0, 255, "000000FF"},
		{"transparent", 255, 255, 255, 0, "FFFFFF00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RGBAToColor(tt.r, tt.g, tt.b, tt.a)
			if got != tt.want {
				t.Errorf("RGBAToColor(%d,%d,%d,%d) = %q, want %q",
					tt.r, tt.g, tt.b, tt.a, got, tt.want)
			}
		})
	}
}

func TestColorToRGB(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
	}{
		{"red", "FF0000FF", "#FF0000"},
		{"green", "00FF00FF", "#00FF00"},
		{"transparent white", "FFFFFF00", "#FFFFFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ColorToRGB(tt.color)
			if err != nil {
				t.Fatalf("ColorToRGB(%q) error = %v", tt.color, err)
			}
			if got != tt.want {
				t.Errorf("ColorToRGB(%q) = %q, want %q", tt.color, got, tt.want)
			}
		})
	}
}

func TestColorWithAlpha(t *testing.T) {
	tests := []struct {
		name  string
		color string
		alpha byte
		want  string
	}{
		{"opaque to transparent", "FF0000FF", 0x00, "FF000000"},
		{"transparent to opaque", "FF000000", 0xFF, "FF0000FF"},
		{"half transparent", "00FF00FF", 0x80, "00FF0080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ColorWithAlpha(tt.color, tt.alpha)
			if err != nil {
				t.Fatalf("ColorWithAlpha(%q, %d) error = %v", tt.color, tt.alpha, err)
			}
			if got != tt.want {
				t.Errorf("ColorWithAlpha(%q, %d) = %q, want %q",
					tt.color, tt.alpha, got, tt.want)
			}
		})
	}
}

func TestColorConstants(t *testing.T) {
	// Verify all constants are valid
	constants := map[string]string{
		"ColorBlack":       ColorBlack,
		"ColorWhite":       ColorWhite,
		"ColorRed":         ColorRed,
		"ColorGreen":       ColorGreen,
		"ColorBlue":        ColorBlue,
		"ColorYellow":      ColorYellow,
		"ColorCyan":        ColorCyan,
		"ColorMagenta":     ColorMagenta,
		"ColorGray":        ColorGray,
		"ColorLightGray":   ColorLightGray,
		"ColorDarkGray":    ColorDarkGray,
		"ColorTransparent": ColorTransparent,
	}

	for name, color := range constants {
		t.Run(name, func(t *testing.T) {
			if err := ValidateColor(color); err != nil {
				t.Errorf("%s = %q is invalid: %v", name, color, err)
			}
		})
	}
}
