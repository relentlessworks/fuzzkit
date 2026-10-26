package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/fuzzkit/internal/model"
)

// Handler holds the HTTP handlers for fuzzkit.
type Handler struct {
	noAuth bool
}

// NewHandler creates a new API handler.
func NewHandler(noAuth bool) *Handler {
	return &Handler{noAuth: noAuth}
}

// Routes returns the HTTP mux with all routes registered.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Help / self-documentation
	mux.HandleFunc("/help", h.handleHelp)
	mux.HandleFunc("/.well-known/agent.md", h.handleHelp)

	// Algorithm listings
	mux.HandleFunc("/algorithms", h.handleAlgorithms)
	mux.HandleFunc("/algorithms/phonetic", h.handlePhoneticAlgorithms)

	// Distance endpoints
	mux.HandleFunc("/levenshtein", h.handleLevenshtein)
	mux.HandleFunc("/damerau", h.handleDamerau)
	mux.HandleFunc("/hamming", h.handleHamming)

	// Similarity endpoints
	mux.HandleFunc("/jaro", h.handleJaro)
	mux.HandleFunc("/jaro-winkler", h.handleJaroWinkler)
	mux.HandleFunc("/ratio", h.handleRatio)
	mux.HandleFunc("/sorensen-dice", h.handleSorensenDice)
	mux.HandleFunc("/ngram", h.handleNGram)
	mux.HandleFunc("/trigram", h.handleTrigram)

	// Compare all
	mux.HandleFunc("/compare", h.handleCompare)

	// Phonetic endpoints
	mux.HandleFunc("/soundex", h.handleSoundex)
	mux.HandleFunc("/metaphone", h.handleMetaphone)
	mux.HandleFunc("/nysiis", h.handleNysiis)
	mux.HandleFunc("/phonetic", h.handlePhonetic)
	mux.HandleFunc("/phonetic/match", h.handlePhoneticMatch)

	// Fuzzy find / filter / rank
	mux.HandleFunc("/find", h.handleFind)
	mux.HandleFunc("/filter", h.handleFilter)
	mux.HandleFunc("/rank", h.handleRank)

	// MCP endpoint
	mux.HandleFunc("/mcp", h.handleMCP)

	// Root
	mux.HandleFunc("/", h.handleRoot)

	return mux
}

// handleRoot provides a simple root response.
func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, r, http.StatusNotFound, "not found", "GET /help for available endpoints")
		return
	}
	writeText(w, http.StatusOK, "fuzzkit — agentic-first fuzzy string matching service\nGET /help for usage\n")
}

// handleHelp returns the self-documenting help text.
func (h *Handler) handleHelp(w http.ResponseWriter, r *http.Request) {
	help := `fuzzkit — Fuzzy String Matching Service
=====================================

Agentic-first fuzzy string matching service. Compute edit distances, similarity
scores, phonetic codes, and find best matches in lists.

AUTH
  This service is stateless. No authentication required for computation endpoints.
  (Auth middleware available for multi-tenant deployments.)

DISTANCE ENDPOINTS
  GET  /levenshtein?s1=STRING&s2=STRING        Levenshtein edit distance
  GET  /damerau?s1=STRING&s2=STRING             Damerau-Levenshtein distance (with transpositions)
  GET  /hamming?s1=STRING&s2=STRING             Hamming distance (equal-length strings only)

SIMILARITY ENDPOINTS (0.0 to 1.0)
  GET  /jaro?s1=STRING&s2=STRING                Jaro similarity
  GET  /jaro-winkler?s1=STRING&s2=STRING        Jaro-Winkler similarity (prefix-weighted)
  GET  /ratio?s1=STRING&s2=STRING               Difflib-style ratio
  GET  /sorensen-dice?s1=STRING&s2=STRING       Sørensen-Dice coefficient (bigram-based)
  GET  /ngram?s1=STRING&s2=STRING&n=2           N-gram similarity (default n=2)
  GET  /trigram?s1=STRING&s2=STRING             Trigram (3-gram) similarity

COMPARE ALL
  GET  /compare?s1=STRING&s2=STRING             All similarity metrics at once

PHONETIC ENDPOINTS
  GET  /soundex?s=STRING                        Soundex code (4 chars: letter + 3 digits)
  GET  /metaphone?s=STRING                      Metaphone code
  GET  /nysiis?s=STRING                         NYSIIS code
  GET  /phonetic?s=STRING&algorithm=soundex     Phonetic code (soundex|metaphone|nysiis)
  GET  /phonetic/match?s1=STRING&s2=STRING&algorithm=soundex  Check phonetic match

FUZZY SEARCH
  GET  /find?q=QUERY&candidates=a,b,c&algorithm=jaro-winkler   Best match from list
  GET  /filter?q=QUERY&candidates=a,b,c&algorithm=jaro-winkler&threshold=0.7  Filtered matches
  GET  /rank?q=QUERY&candidates=a,b,c&algorithm=jaro-winkler   All candidates ranked

ALGORITHMS
  GET  /algorithms                              List all similarity algorithms
  GET  /algorithms/phonetic                     List all phonetic algorithms

MCP
  POST /mcp                                     Model Context Protocol (JSON-RPC 2.0)

RESPONSE FORMAT
  Plain text by default (key=value pairs, one per line).
  JSON via Accept: application/json header or ?format=json query param.

ERRORS
  error: message | hint: what to do next

EXAMPLES
  curl 'http://localhost:8427/levenshtein?s1=kitten&s2=sitting'
  -> distance=3 ratio=0.571

  curl 'http://localhost:8427/jaro-winkler?s1=martha&s2=marhta'
  -> score=0.961

  curl 'http://localhost:8427/find?q=appel&candidates=apple,orange,banana,apricot&algorithm=jaro-winkler'
  -> match=apple score=0.933 algorithm=jaro-winkler

  curl 'http://localhost:8427/soundex?s=Robert'
  -> code=R163
`
	writeText(w, http.StatusOK, help)
}

// handleAlgorithms lists all similarity algorithms.
func (h *Handler) handleAlgorithms(w http.ResponseWriter, r *http.Request) {
	algorithms := model.AlgorithmList()
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"algorithms": algorithms,
			"count":      len(algorithms),
		})
	} else {
		var sb strings.Builder
		for _, a := range algorithms {
			sb.WriteString(a + "\n")
		}
		writeText(w, http.StatusOK, sb.String())
	}
}

// handlePhoneticAlgorithms lists all phonetic algorithms.
func (h *Handler) handlePhoneticAlgorithms(w http.ResponseWriter, r *http.Request) {
	algorithms := model.PhoneticAlgorithmList()
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"algorithms": algorithms,
			"count":      len(algorithms),
		})
	} else {
		var sb strings.Builder
		for _, a := range algorithms {
			sb.WriteString(a + "\n")
		}
		writeText(w, http.StatusOK, sb.String())
	}
}

// getTwoStrings extracts s1 and s2 from query params.
func getTwoStrings(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	s1 := r.URL.Query().Get("s1")
	s2 := r.URL.Query().Get("s2")
	if s1 == "" || s2 == "" {
		writeError(w, r, http.StatusBadRequest,
			"missing required parameters",
			"provide both s1 and s2 query parameters, e.g. ?s1=kitten&s2=sitting")
		return "", "", false
	}
	return s1, s2, true
}

// handleLevenshtein computes Levenshtein distance.
func (h *Handler) handleLevenshtein(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	dist := model.LevenshteinDistance(s1, s2)
	ratio := model.LevenshteinRatio(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"distance": dist,
		"ratio":    fmt.Sprintf("%.4f", ratio),
		"s1":       s1,
		"s2":       s2,
	})
}

// handleDamerau computes Damerau-Levenshtein distance.
func (h *Handler) handleDamerau(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	dist := model.DamerauLevenshteinDistance(s1, s2)
	ratio := model.DamerauLevenshteinRatio(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"distance": dist,
		"ratio":    fmt.Sprintf("%.4f", ratio),
		"s1":       s1,
		"s2":       s2,
	})
}

// handleHamming computes Hamming distance.
func (h *Handler) handleHamming(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	dist := model.HammingDistance(s1, s2)
	if dist < 0 {
		writeError(w, r, http.StatusBadRequest,
			"strings must be of equal length for Hamming distance",
			fmt.Sprintf("s1 has %d characters, s2 has %d characters — use /levenshtein for unequal-length strings", len([]rune(s1)), len([]rune(s2))))
		return
	}
	r1 := []rune(s1)
	ratio := 1.0
	if len(r1) > 0 {
		ratio = 1.0 - float64(dist)/float64(len(r1))
	}
	writeRecord(w, r, map[string]interface{}{
		"distance": dist,
		"ratio":    fmt.Sprintf("%.4f", ratio),
		"s1":       s1,
		"s2":       s2,
	})
}

// handleJaro computes Jaro similarity.
func (h *Handler) handleJaro(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	score := model.JaroSimilarity(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"s1":    s1,
		"s2":    s2,
	})
}

// handleJaroWinkler computes Jaro-Winkler similarity.
func (h *Handler) handleJaroWinkler(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	score := model.JaroWinklerSimilarity(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"s1":    s1,
		"s2":    s2,
	})
}

// handleRatio computes difflib-style ratio.
func (h *Handler) handleRatio(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	score := model.Ratio(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"s1":    s1,
		"s2":    s2,
	})
}

// handleSorensenDice computes Sørensen-Dice coefficient.
func (h *Handler) handleSorensenDice(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	score := model.SorensenDiceCoefficient(s1, s2)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"s1":    s1,
		"s2":    s2,
	})
}

// handleNGram computes n-gram similarity.
func (h *Handler) handleNGram(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	n := 2
	if nStr := r.URL.Query().Get("n"); nStr != "" {
		if parsed, err := parseInt(nStr); err == nil && parsed >= 1 {
			n = parsed
		}
	}
	score := model.NGramSimilarity(s1, s2, n)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"n":     n,
		"s1":    s1,
		"s2":    s2,
	})
}

// handleTrigram computes trigram similarity.
func (h *Handler) handleTrigram(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	score := model.NGramSimilarity(s1, s2, 3)
	writeRecord(w, r, map[string]interface{}{
		"score": fmt.Sprintf("%.4f", score),
		"s1":    s1,
		"s2":    s2,
	})
}

// handleCompare computes all similarity metrics.
func (h *Handler) handleCompare(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	results := model.CompareAll(s1, s2)
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"s1":     s1,
			"s2":     s2,
			"results": results,
		})
	} else {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("s1=%s s2=%s\n", s1, s2))
		for _, res := range results {
			line := fmt.Sprintf("algorithm=%s score=%.4f", res.Algorithm, res.Score)
			if res.Distance > 0 || res.Algorithm == "levenshtein" || res.Algorithm == "damerau-levenshtein" || res.Algorithm == "hamming" {
				line += fmt.Sprintf(" distance=%d", res.Distance)
			}
			sb.WriteString(line + "\n")
		}
		writeText(w, http.StatusOK, sb.String())
	}
}

// handleSoundex computes Soundex code.
func (h *Handler) handleSoundex(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")
	if s == "" {
		writeError(w, r, http.StatusBadRequest, "missing required parameter", "provide s query parameter, e.g. ?s=Robert")
		return
	}
	code := model.Soundex(s)
	writeRecord(w, r, map[string]interface{}{
		"code": code,
		"s":    s,
	})
}

// handleMetaphone computes Metaphone code.
func (h *Handler) handleMetaphone(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")
	if s == "" {
		writeError(w, r, http.StatusBadRequest, "missing required parameter", "provide s query parameter, e.g. ?s=Robert")
		return
	}
	code := model.Metaphone(s)
	writeRecord(w, r, map[string]interface{}{
		"code": code,
		"s":    s,
	})
}

// handleNysiis computes NYSIIS code.
func (h *Handler) handleNysiis(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")
	if s == "" {
		writeError(w, r, http.StatusBadRequest, "missing required parameter", "provide s query parameter, e.g. ?s=Robert")
		return
	}
	code := model.Nysiis(s)
	writeRecord(w, r, map[string]interface{}{
		"code": code,
		"s":    s,
	})
}

// handlePhonetic computes phonetic code with specified algorithm.
func (h *Handler) handlePhonetic(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")
	if s == "" {
		writeError(w, r, http.StatusBadRequest, "missing required parameter", "provide s query parameter, e.g. ?s=Robert&algorithm=soundex")
		return
	}
	algorithm := r.URL.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = "soundex"
	}
	code := model.PhoneticCode(s, algorithm)
	writeRecord(w, r, map[string]interface{}{
		"code":      code,
		"s":         s,
		"algorithm": algorithm,
	})
}

// handlePhoneticMatch checks if two strings have the same phonetic code.
func (h *Handler) handlePhoneticMatch(w http.ResponseWriter, r *http.Request) {
	s1, s2, ok := getTwoStrings(w, r)
	if !ok {
		return
	}
	algorithm := r.URL.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = "soundex"
	}
	match, code := model.PhoneticMatch(s1, s2, algorithm)
	code2 := model.PhoneticCode(s2, algorithm)
	writeRecord(w, r, map[string]interface{}{
		"match":     match,
		"code1":     code,
		"code2":     code2,
		"algorithm": algorithm,
		"s1":        s1,
		"s2":        s2,
	})
}

// handleFind finds the best match from candidates.
func (h *Handler) handleFind(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	candidatesStr := r.URL.Query().Get("candidates")
	if q == "" || candidatesStr == "" {
		writeError(w, r, http.StatusBadRequest,
			"missing required parameters",
			"provide q and candidates query parameters, e.g. ?q=appel&candidates=apple,orange,banana")
		return
	}
	algorithm := r.URL.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = "jaro-winkler"
	}
	candidates := parseCandidates(candidatesStr)
	result, found := model.FuzzyFind(q, candidates, algorithm)
	if !found {
		writeError(w, r, http.StatusNotFound, "no candidates provided", "provide at least one candidate in the candidates parameter")
		return
	}
	writeRecord(w, r, map[string]interface{}{
		"match":     result.String,
		"score":     fmt.Sprintf("%.4f", result.Score),
		"algorithm": algorithm,
		"query":     q,
	})
}

// handleFilter filters candidates by threshold.
func (h *Handler) handleFilter(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	candidatesStr := r.URL.Query().Get("candidates")
	if q == "" || candidatesStr == "" {
		writeError(w, r, http.StatusBadRequest,
			"missing required parameters",
			"provide q and candidates query parameters, e.g. ?q=appel&candidates=apple,orange,banana&threshold=0.7")
		return
	}
	algorithm := r.URL.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = "jaro-winkler"
	}
	threshold := 0.7
	if tStr := r.URL.Query().Get("threshold"); tStr != "" {
		if parsed, err := parseFloat(tStr); err == nil {
			threshold = parsed
		}
	}
	candidates := parseCandidates(candidatesStr)
	results := model.FuzzyFilter(q, candidates, algorithm, threshold)
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"query":      q,
			"algorithm":  algorithm,
			"threshold":  threshold,
			"count":      len(results),
			"results":    results,
		})
	} else {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("query=%s algorithm=%s threshold=%.2f count=%d\n", q, algorithm, threshold, len(results)))
		for _, res := range results {
			sb.WriteString(fmt.Sprintf("match=%s score=%.4f\n", res.String, res.Score))
		}
		writeText(w, http.StatusOK, sb.String())
	}
}

// handleRank ranks all candidates.
func (h *Handler) handleRank(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	candidatesStr := r.URL.Query().Get("candidates")
	if q == "" || candidatesStr == "" {
		writeError(w, r, http.StatusBadRequest,
			"missing required parameters",
			"provide q and candidates query parameters, e.g. ?q=appel&candidates=apple,orange,banana")
		return
	}
	algorithm := r.URL.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = "jaro-winkler"
	}
	candidates := parseCandidates(candidatesStr)
	results := model.FuzzyRank(q, candidates, algorithm)
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"query":     q,
			"algorithm": algorithm,
			"count":     len(results),
			"results":   results,
		})
	} else {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("query=%s algorithm=%s count=%d\n", q, algorithm, len(results)))
		for _, res := range results {
			sb.WriteString(fmt.Sprintf("match=%s score=%.4f\n", res.String, res.Score))
		}
		writeText(w, http.StatusOK, sb.String())
	}
}

// parseCandidates splits a comma-separated string into a list.
func parseCandidates(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// parseInt parses an integer from a string.
func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

// parseFloat parses a float from a string.
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// handleMCP handles Model Context Protocol JSON-RPC 2.0 requests.
func (h *Handler) handleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "POST JSON-RPC 2.0 request to /mcp")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid JSON", "send a valid JSON-RPC 2.0 request body")
		return
	}

	method, _ := req["method"].(string)
	id := req["id"]

	w.Header().Set("Content-Type", "application/json")

	switch method {
	case "initialize":
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "fuzzkit",
					"version": "0.1.0",
				},
			},
		})

	case "tools/list":
		tools := mcpTools()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"tools": tools,
			},
		})

	case "tools/call":
		params, _ := req["params"].(map[string]interface{})
		toolName, _ := params["name"].(string)
		args, _ := params["arguments"].(map[string]interface{})

		result := h.handleMCPToolCall(toolName, args)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": result,
					},
				},
			},
		})

	default:
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": "method not found",
			},
		})
	}
}

func mcpTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "levenshtein",
			"description": "Compute Levenshtein edit distance between two strings",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string", "description": "First string"},
					"s2": map[string]interface{}{"type": "string", "description": "Second string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "damerau_levenshtein",
			"description": "Compute Damerau-Levenshtein distance (with transpositions)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "hamming",
			"description": "Compute Hamming distance (equal-length strings only)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "jaro",
			"description": "Compute Jaro similarity (0.0 to 1.0)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "jaro_winkler",
			"description": "Compute Jaro-Winkler similarity (prefix-weighted, 0.0 to 1.0)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "ratio",
			"description": "Compute difflib-style similarity ratio (0.0 to 1.0)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "sorensen_dice",
			"description": "Compute Sørensen-Dice coefficient (bigram-based, 0.0 to 1.0)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "ngram",
			"description": "Compute N-gram similarity (0.0 to 1.0)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
					"n":  map[string]interface{}{"type": "integer", "description": "N-gram size (default 2)"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "compare",
			"description": "Compute all similarity metrics between two strings",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s1": map[string]interface{}{"type": "string"},
					"s2": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s1", "s2"},
			},
		},
		{
			"name":        "soundex",
			"description": "Compute Soundex phonetic code",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s"},
			},
		},
		{
			"name":        "metaphone",
			"description": "Compute Metaphone phonetic code",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s"},
			},
		},
		{
			"name":        "nysiis",
			"description": "Compute NYSIIS phonetic code",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"s": map[string]interface{}{"type": "string"},
				},
				"required": []string{"s"},
			},
		},
		{
			"name":        "fuzzy_find",
			"description": "Find the best matching string from a list of candidates",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":      map[string]interface{}{"type": "string"},
					"candidates": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"algorithm":  map[string]interface{}{"type": "string", "description": "Algorithm: levenshtein, damerau, jaro, jaro-winkler, ratio, sorensen-dice, ngram, trigram (default: jaro-winkler)"},
				},
				"required": []string{"query", "candidates"},
			},
		},
		{
			"name":        "fuzzy_filter",
			"description": "Filter candidates by similarity threshold",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":      map[string]interface{}{"type": "string"},
					"candidates": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"algorithm":  map[string]interface{}{"type": "string"},
					"threshold":  map[string]interface{}{"type": "number", "description": "Minimum similarity (default 0.7)"},
				},
				"required": []string{"query", "candidates"},
			},
		},
		{
			"name":        "fuzzy_rank",
			"description": "Rank all candidates by similarity to query",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":      map[string]interface{}{"type": "string"},
					"candidates": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"algorithm":  map[string]interface{}{"type": "string"},
				},
				"required": []string{"query", "candidates"},
			},
		},
	}
}

func (h *Handler) handleMCPToolCall(tool string, args map[string]interface{}) string {
	switch tool {
	case "levenshtein":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		dist := model.LevenshteinDistance(s1, s2)
		ratio := model.LevenshteinRatio(s1, s2)
		return fmt.Sprintf("distance=%d ratio=%.4f s1=%s s2=%s", dist, ratio, s1, s2)

	case "damerau_levenshtein":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		dist := model.DamerauLevenshteinDistance(s1, s2)
		ratio := model.DamerauLevenshteinRatio(s1, s2)
		return fmt.Sprintf("distance=%d ratio=%.4f s1=%s s2=%s", dist, ratio, s1, s2)

	case "hamming":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		dist := model.HammingDistance(s1, s2)
		if dist < 0 {
			return fmt.Sprintf("error: strings must be equal length | hint: s1 has %d chars, s2 has %d chars", len([]rune(s1)), len([]rune(s2)))
		}
		return fmt.Sprintf("distance=%d s1=%s s2=%s", dist, s1, s2)

	case "jaro":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		score := model.JaroSimilarity(s1, s2)
		return fmt.Sprintf("score=%.4f s1=%s s2=%s", score, s1, s2)

	case "jaro_winkler":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		score := model.JaroWinklerSimilarity(s1, s2)
		return fmt.Sprintf("score=%.4f s1=%s s2=%s", score, s1, s2)

	case "ratio":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		score := model.Ratio(s1, s2)
		return fmt.Sprintf("score=%.4f s1=%s s2=%s", score, s1, s2)

	case "sorensen_dice":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		score := model.SorensenDiceCoefficient(s1, s2)
		return fmt.Sprintf("score=%.4f s1=%s s2=%s", score, s1, s2)

	case "ngram":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		n := 2
		if nFloat, ok := args["n"].(float64); ok {
			n = int(nFloat)
		}
		score := model.NGramSimilarity(s1, s2, n)
		return fmt.Sprintf("score=%.4f n=%d s1=%s s2=%s", score, n, s1, s2)

	case "compare":
		s1, _ := args["s1"].(string)
		s2, _ := args["s2"].(string)
		results := model.CompareAll(s1, s2)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("s1=%s s2=%s\n", s1, s2))
		for _, res := range results {
			line := fmt.Sprintf("algorithm=%s score=%.4f", res.Algorithm, res.Score)
			if res.Distance > 0 || res.Algorithm == "levenshtein" || res.Algorithm == "damerau-levenshtein" || res.Algorithm == "hamming" {
				line += fmt.Sprintf(" distance=%d", res.Distance)
			}
			sb.WriteString(line + "\n")
		}
		return sb.String()

	case "soundex":
		s, _ := args["s"].(string)
		code := model.Soundex(s)
		return fmt.Sprintf("code=%s s=%s", code, s)

	case "metaphone":
		s, _ := args["s"].(string)
		code := model.Metaphone(s)
		return fmt.Sprintf("code=%s s=%s", code, s)

	case "nysiis":
		s, _ := args["s"].(string)
		code := model.Nysiis(s)
		return fmt.Sprintf("code=%s s=%s", code, s)

	case "fuzzy_find":
		q, _ := args["query"].(string)
		algorithm, _ := args["algorithm"].(string)
		if algorithm == "" {
			algorithm = "jaro-winkler"
		}
		candidates := toStringSlice(args["candidates"])
		result, found := model.FuzzyFind(q, candidates, algorithm)
		if !found {
			return "error: no candidates provided | hint: provide at least one candidate"
		}
		return fmt.Sprintf("match=%s score=%.4f algorithm=%s query=%s", result.String, result.Score, algorithm, q)

	case "fuzzy_filter":
		q, _ := args["query"].(string)
		algorithm, _ := args["algorithm"].(string)
		if algorithm == "" {
			algorithm = "jaro-winkler"
		}
		threshold := 0.7
		if t, ok := args["threshold"].(float64); ok {
			threshold = t
		}
		candidates := toStringSlice(args["candidates"])
		results := model.FuzzyFilter(q, candidates, algorithm, threshold)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("query=%s algorithm=%s threshold=%.2f count=%d\n", q, algorithm, threshold, len(results)))
		for _, res := range results {
			sb.WriteString(fmt.Sprintf("match=%s score=%.4f\n", res.String, res.Score))
		}
		return sb.String()

	case "fuzzy_rank":
		q, _ := args["query"].(string)
		algorithm, _ := args["algorithm"].(string)
		if algorithm == "" {
			algorithm = "jaro-winkler"
		}
		candidates := toStringSlice(args["candidates"])
		results := model.FuzzyRank(q, candidates, algorithm)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("query=%s algorithm=%s count=%d\n", q, algorithm, len(results)))
		for _, res := range results {
			sb.WriteString(fmt.Sprintf("match=%s score=%.4f\n", res.String, res.Score))
		}
		return sb.String()

	default:
		return fmt.Sprintf("error: unknown tool '%s' | hint: use tools/list to see available tools", tool)
	}
}

func toStringSlice(v interface{}) []string {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
