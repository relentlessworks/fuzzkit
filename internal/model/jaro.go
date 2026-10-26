package model

import "math"

// JaroSimilarity computes the Jaro similarity score between two strings.
// Returns a float64 between 0.0 (no similarity) and 1.0 (exact match).
// The Jaro similarity is based on the number and order of matching characters.
func JaroSimilarity(s1, s2 string) float64 {
	r1 := []rune(s1)
	r2 := []rune(s2)

	m, n := len(r1), len(r2)
	if m == 0 && n == 0 {
		return 1.0
	}
	if m == 0 || n == 0 {
		return 0.0
	}

	// Maximum matching distance
	matchDist := max2(m, n)/2 - 1
	if matchDist < 0 {
		matchDist = 0
	}

	s1Matches := make([]bool, m)
	s2Matches := make([]bool, n)

	matches := 0

	// Find matches
	for i := 0; i < m; i++ {
		start := i - matchDist
		if start < 0 {
			start = 0
		}
		end := i + matchDist + 1
		if end > n {
			end = n
		}
		for j := start; j < end; j++ {
			if s2Matches[j] {
				continue
			}
			if r1[i] != r2[j] {
				continue
			}
			s1Matches[i] = true
			s2Matches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	transpositions := 0
	k := 0
	for i := 0; i < m; i++ {
		if !s1Matches[i] {
			continue
		}
		for !s2Matches[k] {
			k++
		}
		if r1[i] != r2[k] {
			transpositions++
		}
		k++
	}
	transpositions /= 2

	return (float64(matches)/float64(m) +
		float64(matches)/float64(n) +
		(float64(matches)-float64(transpositions))/float64(matches)) / 3.0
}

// JaroWinklerSimilarity computes the Jaro-Winkler similarity, which gives
// higher weight to strings that match from the beginning (common prefix).
// The prefix length is capped at 4 characters, and the scaling factor is 0.1.
func JaroWinklerSimilarity(s1, s2 string) float64 {
	jaro := JaroSimilarity(s1, s2)

	// Common prefix length (max 4)
	prefixLen := 0
	r1 := []rune(s1)
	r2 := []rune(s2)
	maxPrefix := min2(min2(len(r1), len(r2)), 4)
	for i := 0; i < maxPrefix; i++ {
		if r1[i] == r2[i] {
			prefixLen++
		} else {
			break
		}
	}

	return jaro + float64(prefixLen)*0.1*(1.0-jaro)
}

// LevenshteinRatio returns a normalized similarity ratio (0.0 to 1.0)
// based on Levenshtein distance.
func LevenshteinRatio(s1, s2 string) float64 {
	dist := LevenshteinDistance(s1, s2)
	maxLen := max2(len([]rune(s1)), len([]rune(s2)))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxLen)
}

// DamerauLevenshteinRatio returns a normalized similarity ratio (0.0 to 1.0)
// based on Damerau-Levenshtein distance.
func DamerauLevenshteinRatio(s1, s2 string) float64 {
	dist := DamerauLevenshteinDistance(s1, s2)
	maxLen := max2(len([]rune(s1)), len([]rune(s2)))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxLen)
}

// Ratio computes a similarity ratio using a difflib-style algorithm.
// Returns a float64 between 0.0 and 1.0.
func Ratio(s1, s2 string) float64 {
	r1 := []rune(s1)
	r2 := []rune(s2)
	m, n := len(r1), len(r2)
	if m == 0 && n == 0 {
		return 1.0
	}
	if m == 0 || n == 0 {
		return 0.0
	}

	// Count total matching characters using recursive longest-match blocks
	matches := countMatchingChars(r1, r2)
	if matches == 0 {
		return 0.0
	}

	return 2.0 * float64(matches) / float64(m+n)
}

// countMatchingChars recursively counts matching characters by finding
// the longest matching block and then recursing on the parts before and after.
func countMatchingChars(r1, r2 []rune) int {
	if len(r1) == 0 || len(r2) == 0 {
		return 0
	}

	match := longestMatch(r1, r2)
	if match.size == 0 {
		return 0
	}

	before := countMatchingChars(r1[:match.s1Start], r2[:match.s2Start])
	after := countMatchingChars(r1[match.s1Start+match.size:], r2[match.s2Start+match.size:])

	return match.size + before + after
}

type matchBlock struct {
	s1Start int
	s2Start int
	size    int
}

func longestMatch(r1, r2 []rune) matchBlock {
	m, n := len(r1), len(r2)
	best := matchBlock{}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			k := 0
			for i+k < m && j+k < n && r1[i+k] == r2[j+k] {
				k++
			}
			if k > best.size {
				best = matchBlock{s1Start: i, s2Start: j, size: k}
			}
		}
	}

	return best
}

// SorensenDiceCoefficient computes the Sørensen-Dice coefficient
// using bigrams (2-grams). Returns a value between 0.0 and 1.0.
func SorensenDiceCoefficient(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if len(s1) < 2 || len(s2) < 2 {
		return 0.0
	}

	b1 := bigrams(s1)
	b2 := bigrams(s2)

	if len(b1) == 0 || len(b2) == 0 {
		return 0.0
	}

	// Count intersection
	intersection := 0
	b2Copy := make(map[string]int)
	for _, b := range b2 {
		b2Copy[b]++
	}
	for _, b := range b1 {
		if b2Copy[b] > 0 {
			intersection++
			b2Copy[b]--
		}
	}

	return 2.0 * float64(intersection) / float64(len(b1)+len(b2))
}

func bigrams(s string) []string {
	r := []rune(s)
	if len(r) < 2 {
		return nil
	}
	result := make([]string, 0, len(r)-1)
	for i := 0; i < len(r)-1; i++ {
		result = append(result, string(r[i:i+2]))
	}
	return result
}

// NGramSimilarity computes n-gram based similarity between two strings.
// Returns a value between 0.0 and 1.0.
func NGramSimilarity(s1, s2 string, n int) float64 {
	if n < 1 {
		n = 2
	}
	if s1 == s2 {
		return 1.0
	}

	g1 := ngrams(s1, n)
	g2 := ngrams(s2, n)

	if len(g1) == 0 && len(g2) == 0 {
		return 1.0
	}
	if len(g1) == 0 || len(g2) == 0 {
		return 0.0
	}

	// Jaccard similarity on n-gram sets
	g2Copy := make(map[string]int)
	for _, g := range g2 {
		g2Copy[g]++
	}
	intersection := 0
	for _, g := range g1 {
		if g2Copy[g] > 0 {
			intersection++
			g2Copy[g]--
		}
	}

	union := len(g1) + len(g2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func ngrams(s string, n int) []string {
	r := []rune(s)
	if len(r) < n {
		return []string{s}
	}
	result := make([]string, 0, len(r)-n+1)
	for i := 0; i <= len(r)-n; i++ {
		result = append(result, string(r[i:i+n]))
	}
	return result
}

// HammingDistance computes the Hamming distance between two strings.
// The strings must be of equal length. Returns -1 if lengths differ.
func HammingDistance(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	if len(r1) != len(r2) {
		return -1
	}
	dist := 0
	for i := range r1 {
		if r1[i] != r2[i] {
			dist++
		}
	}
	return dist
}

// roundFloat rounds a float64 to the specified number of decimal places.
func roundFloat(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}
