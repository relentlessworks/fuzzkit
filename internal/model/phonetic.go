package model

import (
	"strings"
	"unicode"
)

// Soundex computes the Soundex code for a string.
// Soundex is a phonetic algorithm that encodes names by sound.
// The code is a letter followed by three digits.
func Soundex(s string) string {
	if len(s) == 0 {
		return "0000"
	}

	s = strings.ToUpper(s)
	r := []rune(s)

	// Filter to letters only
	var letters []rune
	for _, ch := range r {
		if unicode.IsLetter(ch) {
			letters = append(letters, ch)
		}
	}
	if len(letters) == 0 {
		return "0000"
	}

	// Soundex digit mapping
	soundexMap := map[rune]rune{
		'B': '1', 'F': '1', 'P': '1', 'V': '1',
		'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
		'D': '3', 'T': '3',
		'L': '4',
		'M': '5', 'N': '5',
		'R': '6',
	}

	result := []rune{letters[0]}

	// Handle first letter: if it has a code, remember it to avoid duplicate
	prevCode := soundexMap[letters[0]]

	for i := 1; i < len(letters) && len(result) < 4; i++ {
		ch := letters[i]
		code, hasCode := soundexMap[ch]
		if !hasCode {
			// H and W are "transparent" — they don't reset prevCode.
			// Vowels (A, E, I, O, U, Y) do reset prevCode.
			if ch != 'H' && ch != 'W' {
				prevCode = 0
			}
			continue
		}
		// Skip if same code as previous (including across vowels/H/W/Y)
		if code != prevCode {
			result = append(result, code)
		}
		prevCode = code
	}

	// Pad with zeros
	for len(result) < 4 {
		result = append(result, '0')
	}

	return string(result)
}

// Metaphone computes the Metaphone code for a string.
// Metaphone is a phonetic algorithm that improves on Soundex.
func Metaphone(s string) string {
	if len(s) == 0 {
		return ""
	}

	// Convert to uppercase and filter to letters
	s = strings.ToUpper(s)
	var input []rune
	for _, ch := range s {
		if unicode.IsLetter(ch) {
			input = append(input, ch)
		}
	}
	if len(input) == 0 {
		return ""
	}

	// Handle initial silent letters
	inputStr := string(input)

	// Apply transformations
	result := metaphoneTransform(inputStr)

	return result
}

func metaphoneTransform(s string) string {
	// Convert to rune slice for manipulation
	chars := []rune(s)
	n := len(chars)
	if n == 0 {
		return ""
	}

	var result []rune

	// Drop duplicate adjacent letters (except C)
	i := 0
	for i < n {
		ch := chars[i]

		// Skip duplicate adjacent letters (except C)
		if i > 0 && ch == chars[i-1] && ch != 'C' {
			i++
			continue
		}

		// Transform based on rules
		switch ch {
		case 'A', 'E', 'I', 'O', 'U':
			if i == 0 {
				result = append(result, ch)
			}
			// Vowels after first position are dropped (except initial)

		case 'B':
			// MB -> M (B is silent after M)
			if !(i > 0 && chars[i-1] == 'M') {
				result = append(result, 'B')
			}

		case 'C':
			if i == 0 {
				// CH -> X (as in "character")
				if i+1 < n && chars[i+1] == 'H' {
					result = append(result, 'X')
					i++ // skip H
				} else {
					result = append(result, 'K')
				}
			} else if i+1 < n && chars[i+1] == 'H' {
				// SCH -> SK (S already output, C becomes K)
				if i > 0 && chars[i-1] == 'S' {
					result = append(result, 'K')
				} else {
					result = append(result, 'X')
				}
				i++ // skip H
			} else if i+1 < n && (chars[i+1] == 'I' || chars[i+1] == 'E' || chars[i+1] == 'Y') {
				// CI, CE, CY -> S
				result = append(result, 'S')
			} else {
				result = append(result, 'K')
			}

		case 'D':
			// DG before I/E/Y -> J (as in "edge")
			if i+1 < n && chars[i+1] == 'G' && i+2 < n && (chars[i+2] == 'I' || chars[i+2] == 'E' || chars[i+2] == 'Y') {
				result = append(result, 'J')
				i++ // skip G
			} else {
				result = append(result, 'T')
			}

		case 'F':
			result = append(result, 'F')

		case 'G':
			if i+1 < n && chars[i+1] == 'H' {
				// GH: silent after vowel, otherwise F
				if i > 0 && isVowel(chars[i-1]) {
					// silent
				} else {
					result = append(result, 'F')
				}
				i++ // skip H
			} else if i+1 < n && (chars[i+1] == 'I' || chars[i+1] == 'E' || chars[i+1] == 'Y') {
				// GI, GE, GY -> J (soft G)
				result = append(result, 'J')
			} else {
				result = append(result, 'K')
			}

		case 'H':
			// Silent after vowel, keep after consonant
			if i > 0 && isVowel(chars[i-1]) {
				// silent
			} else if i+1 < n && isVowel(chars[i+1]) {
				result = append(result, 'H')
			}

		case 'J':
			result = append(result, 'J')

		case 'K':
			// Silent after C
			if !(i > 0 && chars[i-1] == 'C') {
				result = append(result, 'K')
			}

		case 'L':
			result = append(result, 'L')

		case 'M':
			result = append(result, 'M')

		case 'N':
			result = append(result, 'N')

		case 'P':
			if i+1 < n && chars[i+1] == 'H' {
				result = append(result, 'F')
				i++ // skip H
			} else {
				result = append(result, 'P')
			}

		case 'Q':
			result = append(result, 'K')

		case 'R':
			result = append(result, 'R')

		case 'S':
			if i+1 < n && chars[i+1] == 'H' {
				result = append(result, 'X')
				i++ // skip H
			} else if i+1 < n && chars[i+1] == 'I' && i+2 < n && (chars[i+2] == 'A' || chars[i+2] == 'O') {
				// SIA, SIO -> X (as in "sion")
				result = append(result, 'X')
			} else {
				result = append(result, 'S')
			}

		case 'T':
			if i+1 < n && chars[i+1] == 'H' {
				result = append(result, '0') // theta
				i++ // skip H
			} else if i+1 < n && chars[i+1] == 'I' && i+2 < n && (chars[i+2] == 'A' || chars[i+2] == 'O') {
				// TIA, TIO -> X
				result = append(result, 'X')
			} else {
				result = append(result, 'T')
			}

		case 'V':
			result = append(result, 'F')

		case 'W', 'Y':
			// Silent unless followed by vowel
			if i+1 < n && isVowel(chars[i+1]) {
				result = append(result, ch)
			}

		case 'X':
			result = append(result, 'K')
			result = append(result, 'S')

		case 'Z':
			result = append(result, 'S')

		default:
			result = append(result, ch)
		}

		i++
	}

	// Deduplicate adjacent output characters
	if len(result) == 0 {
		return ""
	}
	deduped := []rune{result[0]}
	for _, ch := range result[1:] {
		if ch != deduped[len(deduped)-1] {
			deduped = append(deduped, ch)
		}
	}

	return string(deduped)
}

func isVowel(ch rune) bool {
	switch ch {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

// Nysiis computes the NYSIIS (New York State Identification and Intelligence System) code.
// NYSIIS is a phonetic algorithm that improves on Soundex.
func Nysiis(s string) string {
	if len(s) == 0 {
		return ""
	}

	// Convert to uppercase
	s = strings.ToUpper(s)

	// Filter to letters only
	var input []rune
	for _, ch := range s {
		if unicode.IsLetter(ch) {
			input = append(input, ch)
		}
	}
	if len(input) == 0 {
		return ""
	}

	// Step 1: Handle first character translations
	first := input[0]
	switch first {
	case 'M', 'K', 'P':
		input[0] = 'N'
	case 'S', 'Z':
		input[0] = 'S'
	}

	// Step 2: Translate remaining characters
	var result []rune
	result = append(result, input[0])

	i := 1
	for i < len(input) {
		ch := input[i]

		// EV -> AF
		if ch == 'E' && i+1 < len(input) && input[i+1] == 'V' {
			result = append(result, 'A', 'F')
			i += 2
			continue
		}

		// Translate based on character
		switch ch {
		case 'A', 'E', 'I', 'O', 'U':
			result = append(result, 'A')
		case 'Q':
			result = append(result, 'G')
		case 'Y':
			result = append(result, 'A')
		case 'M', 'K', 'P':
			result = append(result, 'N')
		case 'S', 'Z':
			result = append(result, 'S')
		case 'H':
			// Keep H only if preceded or followed by a vowel
			if (i > 0 && isVowel(input[i-1])) || (i+1 < len(input) && isVowel(input[i+1])) {
				result = append(result, 'H')
			}
		case 'W':
			// W -> A if preceded by vowel
			if i > 0 && isVowel(input[i-1]) {
				result = append(result, 'A')
			}
		default:
			result = append(result, ch)
		}
		i++
	}

	// Step 3: Remove trailing 'S' or 'A'
	code := string(result)
	code = strings.TrimSuffix(code, "S")
	code = strings.TrimSuffix(code, "A")

	// Step 4: Replace trailing AY with Y
	if strings.HasSuffix(code, "AY") {
		code = code[:len(code)-2] + "Y"
	}

	// Step 5: Remove duplicate adjacent characters
	if len(code) == 0 {
		return ""
	}
	var final []rune
	final = append(final, rune(code[0]))
	for _, ch := range code[1:] {
		if ch != final[len(final)-1] {
			final = append(final, ch)
		}
	}

	return string(final)
}
