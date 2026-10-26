package model

import (
	"sort"
	"strings"
)

// MatchResult represents a single fuzzy match result.
type MatchResult struct {
	String    string  `json:"string"`
	Score     float64 `json:"score"`
	Algorithm string  `json:"algorithm"`
}

// FuzzyFind finds the best matching string from a list using the specified algorithm.
// Returns the best match and its score.
func FuzzyFind(query string, candidates []string, algorithm string) (MatchResult, bool) {
	if len(candidates) == 0 {
		return MatchResult{}, false
	}

	results := FuzzyRank(query, candidates, algorithm)
	return results[0], true
}

// FuzzyRank ranks all candidates by similarity to the query using the specified algorithm.
// Results are sorted by score descending.
func FuzzyRank(query string, candidates []string, algorithm string) []MatchResult {
	results := make([]MatchResult, 0, len(candidates))
	for _, c := range candidates {
		score := computeScore(query, c, algorithm)
		results = append(results, MatchResult{
			String:    c,
			Score:     score,
			Algorithm: algorithm,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// FuzzyFilter filters candidates by a minimum similarity threshold.
// Results are sorted by score descending.
func FuzzyFilter(query string, candidates []string, algorithm string, threshold float64) []MatchResult {
	results := make([]MatchResult, 0)
	for _, c := range candidates {
		score := computeScore(query, c, algorithm)
		if score >= threshold {
			results = append(results, MatchResult{
				String:    c,
				Score:     score,
				Algorithm: algorithm,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// computeScore computes the similarity score between two strings using the specified algorithm.
func computeScore(s1, s2 string, algorithm string) float64 {
	switch strings.ToLower(algorithm) {
	case "levenshtein", "lev":
		return LevenshteinRatio(s1, s2)
	case "damerau", "damerau-levenshtein", "dl":
		return DamerauLevenshteinRatio(s1, s2)
	case "jaro":
		return JaroSimilarity(s1, s2)
	case "jaro-winkler", "jw", "jarowinkler":
		return JaroWinklerSimilarity(s1, s2)
	case "ratio", "difflib":
		return Ratio(s1, s2)
	case "sorensen", "sorensen-dice", "dice":
		return SorensenDiceCoefficient(s1, s2)
	case "ngram", "n-gram":
		return NGramSimilarity(s1, s2, 2)
	case "trigram":
		return NGramSimilarity(s1, s2, 3)
	default:
		return JaroWinklerSimilarity(s1, s2)
	}
}

// AlgorithmList returns all supported algorithm names.
func AlgorithmList() []string {
	return []string{
		"levenshtein",
		"damerau",
		"jaro",
		"jaro-winkler",
		"ratio",
		"sorensen-dice",
		"ngram",
		"trigram",
	}
}

// PhoneticAlgorithmList returns all supported phonetic algorithm names.
func PhoneticAlgorithmList() []string {
	return []string{
		"soundex",
		"metaphone",
		"nysiis",
	}
}

// AllAlgorithms returns all supported algorithm names (similarity + phonetic).
func AllAlgorithms() []string {
	return append(AlgorithmList(), PhoneticAlgorithmList()...)
}

// PhoneticMatch checks if two strings have the same phonetic code.
func PhoneticMatch(s1, s2 string, algorithm string) (bool, string) {
	var code1, code2 string
	switch strings.ToLower(algorithm) {
	case "soundex":
		code1 = Soundex(s1)
		code2 = Soundex(s2)
	case "metaphone":
		code1 = Metaphone(s1)
		code2 = Metaphone(s2)
	case "nysiis":
		code1 = Nysiis(s1)
		code2 = Nysiis(s2)
	default:
		code1 = Soundex(s1)
		code2 = Soundex(s2)
	}
	return code1 == code2, code1
}

// PhoneticCode computes the phonetic code for a string using the specified algorithm.
func PhoneticCode(s string, algorithm string) string {
	switch strings.ToLower(algorithm) {
	case "soundex":
		return Soundex(s)
	case "metaphone":
		return Metaphone(s)
	case "nysiis":
		return Nysiis(s)
	default:
		return Soundex(s)
	}
}

// CompareAll computes all similarity scores between two strings.
type CompareResult struct {
	Algorithm string  `json:"algorithm"`
	Score     float64 `json:"score"`
	Distance  int     `json:"distance,omitempty"`
}

// CompareAll computes all similarity metrics between two strings.
func CompareAll(s1, s2 string) []CompareResult {
	results := []CompareResult{
		{Algorithm: "levenshtein", Score: roundFloat(LevenshteinRatio(s1, s2), 4), Distance: LevenshteinDistance(s1, s2)},
		{Algorithm: "damerau-levenshtein", Score: roundFloat(DamerauLevenshteinRatio(s1, s2), 4), Distance: DamerauLevenshteinDistance(s1, s2)},
		{Algorithm: "jaro", Score: roundFloat(JaroSimilarity(s1, s2), 4)},
		{Algorithm: "jaro-winkler", Score: roundFloat(JaroWinklerSimilarity(s1, s2), 4)},
		{Algorithm: "ratio", Score: roundFloat(Ratio(s1, s2), 4)},
		{Algorithm: "sorensen-dice", Score: roundFloat(SorensenDiceCoefficient(s1, s2), 4)},
		{Algorithm: "ngram", Score: roundFloat(NGramSimilarity(s1, s2, 2), 4)},
		{Algorithm: "trigram", Score: roundFloat(NGramSimilarity(s1, s2, 3), 4)},
		{Algorithm: "hamming", Score: hammingRatio(s1, s2), Distance: HammingDistance(s1, s2)},
	}
	return results
}

func hammingRatio(s1, s2 string) float64 {
	dist := HammingDistance(s1, s2)
	if dist < 0 {
		return 0.0
	}
	r1 := []rune(s1)
	if len(r1) == 0 {
		return 1.0
	}
	return roundFloat(1.0-float64(dist)/float64(len(r1)), 4)
}

// CompareAllPhonetic computes all phonetic codes for a string.
type PhoneticResult struct {
	Algorithm string `json:"algorithm"`
	Code      string `json:"code"`
}

// CompareAllPhonetic computes all phonetic codes for a string.
func CompareAllPhonetic(s string) []PhoneticResult {
	return []PhoneticResult{
		{Algorithm: "soundex", Code: Soundex(s)},
		{Algorithm: "metaphone", Code: Metaphone(s)},
		{Algorithm: "nysiis", Code: Nysiis(s)},
	}
}
