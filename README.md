# fuzzkit

> Agentic-first fuzzy string matching service. The agent IS the interface.

Fuzzkit provides fuzzy string matching algorithms over plain HTTP. Compute edit distances, similarity scores, phonetic codes, and find best matches in lists of candidates.

## Quick Start

```bash
# Build and run
make build
./fuzzkit -no-auth

# Or with Go directly
go run ./cmd/fuzzkit -no-auth

# Query
curl 'http://localhost:8427/levenshtein?s1=kitten&s2=sitting'
# distance=3 ratio=0.5714 s1=kitten s2=sitting

curl 'http://localhost:8427/find?q=appel&candidates=apple,orange,banana,apricot&algorithm=jaro-winkler'
# match=apple score=0.9333 algorithm=jaro-winkler query=appel
```

## Principles

- **Agent IS the interface** — No UI, no SDK. The API is the product.
- **Plain text by default** — One labeled, grepable line per record. JSON on demand via `Accept: application/json` or `?format=json`.
- **Instructive errors** — Every 4xx includes a hint telling the agent what to do next.
- **Self-documenting** — `GET /help` returns a one-page operating manual.
- **Single static binary** — Go, CGO_ENABLED=0, zero external dependencies.
- **Zero config defaults** — Runs out of the box.
- **MCP connector** — Speaks Model Context Protocol at `POST /mcp`.

## API Reference

### Distance Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/levenshtein?s1=&s2=` | Levenshtein edit distance |
| GET | `/damerau?s1=&s2=` | Damerau-Levenshtein distance (with transpositions) |
| GET | `/hamming?s1=&s2=` | Hamming distance (equal-length strings only) |

### Similarity Endpoints (0.0 to 1.0)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/jaro?s1=&s2=` | Jaro similarity |
| GET | `/jaro-winkler?s1=&s2=` | Jaro-Winkler similarity (prefix-weighted) |
| GET | `/ratio?s1=&s2=` | Difflib-style ratio |
| GET | `/sorensen-dice?s1=&s2=` | Sørensen-Dice coefficient (bigram-based) |
| GET | `/ngram?s1=&s2=&n=2` | N-gram similarity |
| GET | `/trigram?s1=&s2=` | Trigram (3-gram) similarity |

### Compare All

| Method | Path | Description |
|--------|------|-------------|
| GET | `/compare?s1=&s2=` | All similarity metrics at once |

### Phonetic Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/soundex?s=` | Soundex code |
| GET | `/metaphone?s=` | Metaphone code |
| GET | `/nysiis?s=` | NYSIIS code |
| GET | `/phonetic?s=&algorithm=soundex` | Phonetic code (soundex\|metaphone\|nysiis) |
| GET | `/phonetic/match?s1=&s2=&algorithm=soundex` | Check phonetic match |

### Fuzzy Search

| Method | Path | Description |
|--------|------|-------------|
| GET | `/find?q=&candidates=a,b,c&algorithm=jaro-winkler` | Best match from list |
| GET | `/filter?q=&candidates=a,b,c&algorithm=jaro-winkler&threshold=0.7` | Filtered matches |
| GET | `/rank?q=&candidates=a,b,c&algorithm=jaro-winkler` | All candidates ranked |

### Other

| Method | Path | Description |
|--------|------|-------------|
| GET | `/algorithms` | List all similarity algorithms |
| GET | `/algorithms/phonetic` | List all phonetic algorithms |
| GET | `/help` | Self-documenting help |
| POST | `/mcp` | MCP JSON-RPC 2.0 endpoint |

## Algorithms

### Similarity Algorithms
- **levenshtein** — Edit distance ratio
- **damerau** — Damerau-Levenshtein (with transpositions)
- **jaro** — Jaro similarity
- **jaro-winkler** — Jaro-Winkler (prefix-weighted, default)
- **ratio** — Difflib-style ratio
- **sorensen-dice** — Sørensen-Dice coefficient (bigram-based)
- **ngram** — N-gram similarity (configurable n)
- **trigram** — Trigram (3-gram) similarity

### Phonetic Algorithms
- **soundex** — Soundex (4-char code)
- **metaphone** — Metaphone
- **nysiis** — NYSIIS

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `FUZZKIT_ADDR` | `:8427` | Listen address |
| `-no-auth` | `FUZZKIT_NO_AUTH` | `false` | Disable auth (dev mode) |

## Build

```bash
make build    # CGO_ENABLED=0 go build -trimpath
make test     # go test -race ./...
make vet      # go vet ./...
make run      # build and run with -no-auth
```

## License

MIT
