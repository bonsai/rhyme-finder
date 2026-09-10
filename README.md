# rhyme-finder

Go-first, API-first Japanese rhyme candidate engine.

`rhyme-finder-java` is treated as the reference implementation for rhyme concepts and search behavior; this repository is a clean Go implementation optimized for API serving, LLM/agent integration, and low-latency deterministic ranking.

## Goals

- Fast deterministic candidate ranking.
- OpenAPI + JSON as the stable integration boundary.
- Aggressive mode: prefer useful/creative candidates over overly strict phonetic filtering.
- Strict mode: conservative phonetic candidates.
- Keep the core independent of an LLM.
- Make MeCab/mora analysis replaceable through a future analyzer interface.

## API

```http
POST /v1/rhyme
Content-Type: application/json

{"text":"強引","mode":"aggressive","limit":20}
```

The response contains `score`, `phonetic`, `rhythm`, and `position` components so an LLM can explain or re-rank results without having to reproduce the phonetic engine.

See [`openapi.yaml`](./openapi.yaml).

## Run

```bash
go run ./cmd/rhyme-finder -addr :8080
```

A newline-separated dictionary can be supplied with `-dict`.

## Architecture

```text
LLM / Agent
    ↓ JSON
OpenAPI
    ↓
rhyme-finder (Go)
    ├─ analysis
    ├─ phonetic scoring
    ├─ rhythm scoring
    ├─ aggressive expansion
    └─ ranking
         ↓
      JSON result
```

## Reference

The Java repository remains the algorithm/reference source and is not a port target. New Go behavior should be backed by focused tests and, where appropriate, comparison fixtures against the Java implementation.
