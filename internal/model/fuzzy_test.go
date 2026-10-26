package model

import (
	"testing"
)

func TestFuzzyFind(t *testing.T) {
	candidates := []string{"apple", "orange", "banana", "apricot", "grape"}
	result, found := FuzzyFind("appel", candidates, "jaro-winkler")
	if !found {
		t.Fatal("expected to find a match")
	}
	if result.String != "apple" {
		t.Errorf("expected best match 'apple', got %q (score=%.4f)", result.String, result.Score)
	}
}

func TestFuzzyFilter(t *testing.T) {
	candidates := []string{"apple", "orange", "banana", "apricot", "grape"}
	results := FuzzyFilter("appel", candidates, "jaro-winkler", 0.7)
	if len(results) == 0 {
		t.Fatal("expected at least one result above threshold")
	}
	if results[0].String != "apple" {
		t.Errorf("expected best match 'apple', got %q", results[0].String)
	}
}

func TestFuzzyRank(t *testing.T) {
	candidates := []string{"apple", "orange", "banana", "apricot"}
	results := FuzzyRank("appel", candidates, "jaro-winkler")
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	// Results should be sorted by score descending
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted by score descending: [%d].Score=%.4f > [%d].Score=%.4f", i, results[i].Score, i-1, results[i-1].Score)
		}
	}
}

func TestFuzzyFindLevenshtein(t *testing.T) {
	candidates := []string{"kitten", "sitting", "knitting", "mitten"}
	result, found := FuzzyFind("kitten", candidates, "levenshtein")
	if !found {
		t.Fatal("expected to find a match")
	}
	if result.String != "kitten" {
		t.Errorf("expected exact match 'kitten', got %q", result.String)
	}
}

func TestCompareAll(t *testing.T) {
	results := CompareAll("kitten", "sitting")
	if len(results) == 0 {
		t.Fatal("expected results from CompareAll")
	}
	// Check that levenshtein distance is 3
	for _, r := range results {
		if r.Algorithm == "levenshtein" && r.Distance != 3 {
			t.Errorf("expected levenshtein distance 3, got %d", r.Distance)
		}
	}
}

func TestPhoneticMatch(t *testing.T) {
	match, _ := PhoneticMatch("Robert", "Rupert", "soundex")
	if !match {
		t.Error("expected Robert and Rupert to have the same Soundex code")
	}

	match2, _ := PhoneticMatch("Robert", "Smith", "soundex")
	if match2 {
		t.Error("expected Robert and Smith to have different Soundex codes")
	}
}

func TestPhoneticCode(t *testing.T) {
	code := PhoneticCode("Robert", "soundex")
	if code != "R163" {
		t.Errorf("expected R163, got %s", code)
	}

	code = PhoneticCode("Robert", "metaphone")
	if code != "RBRT" {
		t.Errorf("expected RBRT, got %s", code)
	}
}

func TestAlgorithmList(t *testing.T) {
	algorithms := AlgorithmList()
	if len(algorithms) != 8 {
		t.Errorf("expected 8 algorithms, got %d", len(algorithms))
	}
}

func TestPhoneticAlgorithmList(t *testing.T) {
	algorithms := PhoneticAlgorithmList()
	if len(algorithms) != 3 {
		t.Errorf("expected 3 phonetic algorithms, got %d", len(algorithms))
	}
}
