package model

import "testing"

func TestSoundex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Robert", "R163"},
		{"Rupert", "R163"},
		{"Ashcraft", "A261"},
		{"Tymczak", "T522"},
		{"Pfister", "P236"},
		{"Honeyman", "H555"},
		{"", "0000"},
		{"Lloyd", "L300"},
		{"Washington", "W252"},
	}

	for _, tt := range tests {
		got := Soundex(tt.input)
		if got != tt.want {
			t.Errorf("Soundex(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSoundexHomophones(t *testing.T) {
	// Names that sound the same should have the same Soundex code
	pairs := []struct {
		s1, s2 string
	}{
		{"Robert", "Rupert"},
		{"Schmidt", "Smith"},
		{"Jon", "John"},
	}

	for _, p := range pairs {
		c1 := Soundex(p.s1)
		c2 := Soundex(p.s2)
		if c1 != c2 {
			t.Errorf("Soundex(%q) = %q != Soundex(%q) = %q (expected homophones to match)", p.s1, c1, p.s2, c2)
		}
	}
}

func TestMetaphone(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Robert", "RBRT"},
		{"Smith", "SM0"},
		{"Schmidt", "SKMT"},
		{"example", "EKSMPL"},
		{"phone", "FN"},
		{"dumb", "TM"},
	}

	for _, tt := range tests {
		got := Metaphone(tt.input)
		if got != tt.want {
			t.Errorf("Metaphone(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNysiis(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Robert", "RABART"},
		{"Smith", "SNAT"},
		{"Schmidt", "SCNADT"},
	}

	for _, tt := range tests {
		got := Nysiis(tt.input)
		if got != tt.want {
			t.Errorf("Nysiis(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
