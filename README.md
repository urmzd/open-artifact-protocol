<p align="center">
  <h1 align="center">Open Artifact Protocol</h1>
  <p align="center">
    An open, provider-neutral protocol for LLM artifact updates.
    <br /><br />
    <a href="https://pkg.go.dev/github.com/urmzd/generative-artifact-protocol">Go Docs</a>
    &middot;
    <a href="https://github.com/urmzd/open-artifact-protocol/issues">Report Bug</a>
    &middot;
    <a href="spec/gap.md">Specification</a>
  </p>
</p>

<p align="center">
  <a href="https://github.com/urmzd/open-artifact-protocol/actions/workflows/ci.yml"><img src="https://github.com/urmzd/open-artifact-protocol/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  &nbsp;
  <a href="LICENSE"><img src="https://img.shields.io/github/license/urmzd/open-artifact-protocol" alt="License"></a>
</p>

> **Warning**: This project is `v0`: the protocol, schemas, and APIs are subject to breaking changes without notice until a formal release.

## Status: archived

This repository is archived as of 2026-10-10. The specification, the Go apply engine, and the eval corpus remain available and usable under their existing licenses, but the project is no longer under active development.

**Naming.** The project was renamed from Generative Artifact Protocol to Open Artifact Protocol on 2026-10-10 to reflect what distinguishes it: an open, provider-neutral spec. The wire format is unchanged, so the protocol string (`gap/0.1`), the `<gap:target>` markers, the Go package name (`gap`), and the Go module path (`github.com/urmzd/generative-artifact-protocol`) all keep the original name. The spec and this README still refer to the protocol as GAP.

**Why archived.** Since GAP was published, the industry has converged on a closely similar mechanism for artifact updates.

| Date | Event |
|---|---|
| 2026-02-26 | First GAP commit |
| 2026-03-29 | GAP `v0.1.0` |
| 2026-06-18 | Claude Code Artifacts launches in beta |

Claude Code Artifacts, like GAP, keeps the artifact body in a store, publishes it by reference instead of re-emitting it, applies small targeted edits through the harness, and versions each publish. This note claims convergence on a similar design, not derivation.

**What still differs.**

- **Addressing.** GAP edits name a target by ID (`<gap:target>`) or JSON Pointer. String-replace edits quote the existing text. In two informal edits observed on 2026-10-10, the quoted text was a small share of edit output (roughly 2 to 3 percent), so the token margin between the two approaches is small for small edits.
- **Openness.** GAP is a provider-neutral spec with a small standalone engine. It runs with any model and needs no vendor harness.

**How to read the committed measurements.** The results under [`assets/evals/`](assets/evals/experiments/EXPERIMENT.md) compare GAP against full regeneration of the artifact. They do not compare GAP against string-replace editing, which is now the common baseline. Treat the savings figures as GAP versus full regeneration only.

Open Artifact Protocol, formerly Generative Artifact Protocol (GAP), lets LLMs update structured text artifacts by emitting small JSON envelopes instead of regenerating full files. The native implementation is now Go-only: protocol types, `<gap:target>` marker resolution, JSON Pointer edits, a stateless apply engine, and a versioned in-memory store live in the root Go package.

## Features

- **Envelope system**: `synthesize`, `edit`, and `handle` operations
- **Stateless apply engine**: `Apply(base, envelope) -> artifact, handle`
- **ID-based targeting**: `<gap:target id="ID">` markers for text formats
- **JSON Pointer targeting**: pointer edits for `application/json`
- **Format-agnostic model**: HTML, Markdown, source code, JSON, YAML, SVG, SQL, and other text artifacts
- **Optimistic versioning**: edit envelopes apply only against the expected previous version
- **SSE transport binding**: wire format for streaming with reconnection support in [GAP-SSE](spec/gap-sse.md)

## Installation

```sh
go get github.com/urmzd/generative-artifact-protocol
```

From source:

```sh
git clone https://github.com/urmzd/open-artifact-protocol
cd open-artifact-protocol
just test
```

## Quick Start

```go
package main

import (
	"encoding/json"
	"fmt"

	gap "github.com/urmzd/generative-artifact-protocol"
)

func main() {
	format := "text/html"
	body := `<gap:target id="msg">hello</gap:target>`
	raw, _ := json.Marshal(gap.SynthesizeContentItem{Body: body})

	envelope := gap.Envelope{
		Protocol: gap.ProtocolVersion,
		ID:       "greeting",
		Version:  1,
		Name:     gap.NameSynthesize,
		Meta:     gap.Meta{Format: &format},
		Content:  []json.RawMessage{raw},
	}

	artifact, handle, err := gap.Apply(nil, envelope)
	if err != nil {
		panic(err)
	}
	fmt.Println(artifact.Body)
	fmt.Println(handle.Name)
}
```

## Apply Engine

```go
artifact, handle, err := gap.Apply(baseArtifact, envelope)
```

| Envelope | Direction | Description |
|---|---|---|
| `synthesize` | input | Full artifact content, usually the initial version |
| `edit` | input | Targeted changes via marker ID or JSON Pointer |
| `handle` | output | Lightweight reference returned after every mutation |

The apply engine performs no network or filesystem I/O. It returns errors for malformed envelopes, missing targets, invalid JSON pointers, out-of-bounds array indexes, and unsupported input operation names.

## Store

Use `ArtifactStore` when you want version checks and bounded history:

```go
store := gap.NewArtifactStore(10)
artifact, handle, err := store.Apply(envelope)
```

For edit envelopes, the store enforces `stored_version == envelope.version - 1`. A `synthesize` envelope can create or reset an artifact chain.

## Commands

| Recipe | Description |
|---|---|
| `just build` | Compile the Go package |
| `just test` | Run Go tests |
| `just evalset` | Regenerate SAIGE observations from the GAP eval corpus |
| `just report` | Regenerate eval economics and reliability report from `metrics.json` files |
| `just run` | Run live OpenAI-compatible eval experiments |
| `just check` | Run formatting, vet, and tests |

Equivalent direct commands:

```sh
go test ./...
go vet ./...
go run ./internal/evalsetgen
go run ./internal/evalreport
test -z "$(gofmt -l *.go cmd/gap-eval/*.go evalset/*.go internal/evalsetgen/*.go internal/evalreport/*.go internal/liveeval/*.go)"
```

## Evaluation Corpus

The eval corpus remains under [`assets/evals/`](assets/evals/). The runnable replacement is a Go live runner in [`cmd/gap-eval`](cmd/gap-eval/) plus the SAIGE-compatible eval set in [`assets/evals/saige/`](assets/evals/saige/), generated by [`evalset`](evalset/) as `github.com/urmzd/saige/eval` observations.

Read the benchmark numbers as a risk profile, not a blanket claim that GAP always wins. Raw `comparison` savings measure the protocol attempt itself. New live runs also report `reliability` and `economics` blocks: miss counts/rates, fallback-adjusted savings, retry tax, and init-inclusive amortized savings. If a GAP edit misses, production systems should expect to pay for the failed envelope attempt and then fall back to full regeneration; if an edit applies but corrupts content, only correctness oracles catch it.

The current `gpt-4o-mini` metrics set has 96 metrics files, of which 67 are comparable non-degenerate base-vs-GAP economics runs. In that comparable set, fallback-adjusted edit-token savings are 37.8%, init-inclusive fallback savings are 34.0%, and the miss rate is 14.9%. The strongest savings appear in artifacts above roughly 2 KB, repeated records/pages, dashboards, feeds, catalogs, API payloads, and code files with stable section boundaries. Tiny configs and one-off edits can lose because marker, supervisor inventory, and envelope overhead exceed the full-regeneration baseline.

The runner now treats target discovery as part of the supervisor loop. For non-JSON artifacts it inventories valid `<gap:target>` IDs with `list_targets()`. For JSON artifacts it inventories valid JSON Pointers with `list_paths()`. Generated envelopes are validated before apply; invented targets or invalid paths are rejected and repaired before they can become apply misses.

The previous Rust eval binary and benchmark harness have been removed in favor of the Go runner and SAIGE eval integration. Committed `metrics.json` and `results.md` files are measured artifacts: do not hand-edit measured numbers. Regenerate the SAIGE observation set with `just evalset` after changing experiment prompts or metadata, and regenerate measured result tables with `just report` after changing committed metrics.

## License

This project is dual-licensed:

- **Code** (Go files, build files): [Apache License 2.0](LICENSE)
- **Specification, eval assets, and documentation**: [CC-BY 4.0](LICENSE-CC-BY-4.0)

See [NOTICE](NOTICE) for details. Attribution is required under both licenses.
