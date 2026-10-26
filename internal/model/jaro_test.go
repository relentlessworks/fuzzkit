package model

import (
	"math"
	"testing"
)

func TestJaroSimilarity(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   float64
	}{
		{"MARTHA", "MARHTA", 0.944},
		{"DIXON", "DICKSONX", 0.767},
		{"JELLYFISH", "SMELLYFISH", 0.896},
		{"", "", 1.0},
		{"abc", "", 0.0},
		{"abc", "abc", 1.0},
	}

	for _, tt := range tests {
		got := JaroSimilarity(tt.s1, tt.s2)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("JaroSimilarity(%q, %q) = %.4f, want %.4f", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestJaroWinklerSimilarity(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   float64
	}{
		{"MARTHA", "MARHTA", 0.961},
		{"DIXON", "DICKSONX", 0.813},
		{"", "", 1.0},
		{"abc", "abc", 1.0},
	}

	for _, tt := range tests {
		got := JaroWinklerSimilarity(tt.s1, tt.s2)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("JaroWinklerSimilarity(%q, %q) = %.4f, want %.4f", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   float64
	}{
		{"abcd", "abcd", 1.0},
		{"", "", 1.0},
		{"abc", "", 0.0},
		{"hello world", "world hello", 0.4545},
	}

	for _, tt := range tests {
		got := Ratio(tt.s1, tt.s2)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("Ratio(%q, %q) = %.4f, want %.4f", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestSorensenDiceCoefficient(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   float64
	}{
		{"night", "nacht", 0.25},
		{"abc", "abc", 1.0},
		{"", "", 1.0},
		{"abc", "xyz", 0.0},
		{"hello", "hello", 1.0},
	}

	for _, tt := range tests {
		got := SorensenDiceCoefficient(tt.s1, tt.s2)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("SorensenDiceCoefficient(%q, %q) = %.4f, want %.4f", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestNGramSimilarity(t *testing.T) {
	tests := []struct {
		s1, s2 string
		n      int
		want   float64
	}{
		{"abc", "abc", 2, 1.0},
		{"abc", "abd", 2, 1.0 / 3.0},
		{"", "", 2, 1.0},
		{"abc", "xyz", 2, 0.0},
	}

	for _, tt := range tests {
		got := NGramSimilarity(tt.s1, tt.s2, tt.n)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("NGramSimilarity(%q, %q, %d) = %.4f, want %.4f", tt.s1, tt.s2, tt.n, got, tt.want)
		}
	}
}
