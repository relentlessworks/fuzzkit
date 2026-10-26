package model

import "testing"

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"kitten", "sitting", 3},
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"flaw", "lawn", 2},
		{"gumbo", "gambol", 2},
		{"intention", "execution", 5},
		{"sunday", "saturday", 3},
	}

	for _, tt := range tests {
		got := LevenshteinDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestDamerauLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"ca", "ac", 1},   // transposition
		{"abc", "acb", 1}, // transposition
		{"kitten", "sitting", 3},
		{"", "", 0},
		{"abc", "abc", 0},
		{"abcd", "acbd", 1}, // transposition of b and c
	}

	for _, tt := range tests {
		got := DamerauLevenshteinDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("DamerauLevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestLevenshteinRatio(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   float64
	}{
		{"abc", "abc", 1.0},
		{"", "", 1.0},
		{"abc", "", 0.0},
		{"kitten", "sitting", 4.0 / 7.0},
	}

	for _, tt := range tests {
		got := LevenshteinRatio(tt.s1, tt.s2)
		if absFloat(got-tt.want) > 0.001 {
			t.Errorf("LevenshteinRatio(%q, %q) = %.4f, want %.4f", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestHammingDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"1011101", "1001001", 2},
		{"karolin", "kathrin", 3},
		{"abc", "abcd", -1}, // unequal length
	}

	for _, tt := range tests {
		got := HammingDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("HammingDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
