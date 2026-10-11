<p align="center">
  <h1 align="center">Open Artifact Protocol</h1>
  <p align="center">
    An open, provider-neutral protocol for LLM artifact updates.
    <br /><br />
    <a href="https://pkg.go.dev/github.com/urmzd/open-artifact-protocol">Go Docs</a>
    &middot;
    <a href="https://github.com/urmzd/open-artifact-protocol/issues">Report Bug</a>
    &middot;
    <a href="spec/oap.md">Specification</a>
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

**Naming.** The project was renamed from Generative Artifact Protocol (GAP) to Open Artifact Protocol (OAP) on 2026-10-10 to reflect what distinguishes it: an open, provider-neutral spec. The rename is complete and breaking: the protocol string is now `oap/0.1`, markers are `<oap:target>`, the Go package is `oap`, and the module path is `github.com/urmzd/open-artifact-protocol`. Releases up to `v0.18.0` use the former `gap` identifiers and the former module path. `CHANGELOG.md` keeps the original names as history.

The eval corpus was renamed mechanically along with the code, including recorded model outputs, which originally contained `<gap:target>` and `gap/0.1`. The two names are the same length, so recorded byte counts are unaffected. Token counts were measured with the original names.

**Why archived.** Since the protocol was published as GAP, the industry has converged on a closely similar mechanism for artifact updates.

| Date | Event |
|---|---|
| 2026-02-26 | First commit, as GAP |
| 2026-03-29 | `v0.1.0`, as GAP |
| 2026-06-18 | Claude Code Artifacts launches in beta |

Claude Code Artifacts, like OAP, keeps the artifact body in a store, publishes it by reference instead of re-emitting it, applies small targeted edits through the harness, and versions each publish. This note claims convergence on a similar design, not derivation.

**What still differs.**

- **Addressing.** OAP edits name a target by ID (`<oap:target>`) or JSON Pointer. String-replace edits quote the existing text. In two informal edits observed on 2026-10-10, the quoted text was a small share of edit output (roughly 2 to 3 percent), so the token margin between the two approaches is small for small edits.
- **Openness.** OAP is a provider-neutral spec with a small standalone engine. It runs with any model and needs no vendor harness.

**How to read the committed measurements.** The results under [`assets/evals/`](assets/evals/experiments/EXPERIMENT.md) compare OAP against full regeneration of the artifact. They do not compare OAP against string-replace editing, which is now the common baseline. Treat the savings figures as OAP versus full regeneration only.

### Design notes: what the target IDs were for

These notes record the intent behind the design and what was and was not tested. They are context for readers, not measured results.

**Targets mark data points.** A `<oap:target>` ID names one value or region in the artifact: a metric, a price, a table row, a section. The IDs were meant as stable handles on the data, so that a producer, a reviewer, or another program can look a value up, change it, and check it by name. Cheaper edits were one use of that, not the whole purpose.

**Lookup, iteration, and validation, not only recall.** A string-replace edit is a recall task: the model must reproduce the existing text exactly, and the match fails on any slip. An ID edit is a lookup: the model picks a short name from a list it was just shown. When this protocol was designed, verbatim recall over long artifacts was a real weakness in the models available. Current models with long context have largely outgrown it, and string-replace editing now works reliably in practice. The recall argument for IDs is therefore much weaker than it was.

**Where IDs may still help.** These are possibilities, not findings.

- **Validation.** A check can be written against a named data point ("`total-revenue` equals the sum of the rows") instead of against a text pattern.
- **Classifiers and references.** A classifier, a reviewer, or a second model can point at a target by ID when it flags a problem, and the producer can act on exactly that region.
- **Goals and self-correction.** A goal can be stated as the expected state of named targets. The producer can then compare the artifact to the goal target by target and correct itself, with no need to re-read and re-judge the whole artifact.

**The split this was built for.** OAP separates two roles: an orchestrator that holds only a handle (ID, version, target list), and a maintain context that holds the artifact and edits it. IDs are what make that split possible, because an agent that never saw the body can still address a target by name. String-replace editing cannot do that, since quoting text requires seeing it. Claude Code Artifacts made the opposite choice: one agent writes the file, keeps it in its own context, and edits it by quoting.

**What the evals did not test.** The live OAP flow in this repository sent the maintain context the full artifact on every edit, together with the target list. With the whole body in view, IDs offer little over quoting, which is consistent with the small margins observed. Editing from the handle and target list alone, where input scales with the number of targets instead of the size of the artifact, was never run.

**An open hypothesis.** On small open models, ID addressing may reach a higher exact-match rate than string-replace editing, with the difference closing as models get stronger. This has not been measured here.

Open Artifact Protocol (OAP), formerly Generative Artifact Protocol (GAP), lets LLMs update structured text artifacts by emitting small JSON envelopes instead of regenerating full files. The native implementation is now Go-only: protocol types, `<oap:target>` marker resolution, JSON Pointer edits, a stateless apply engine, and a versioned in-memory store live in the root Go package.

## Features

- **Envelope system**: `synthesize`, `edit`, and `handle` operations
- **Stateless apply engine**: `Apply(base, envelope) -> artifact, handle`
- **ID-based targeting**: `<oap:target id="ID">` markers for text formats
- **JSON Pointer targeting**: pointer edits for `application/json`
- **Format-agnostic model**: HTML, Markdown, source code, JSON, YAML, SVG, SQL, and other text artifacts
- **Optimistic versioning**: edit envelopes apply only against the expected previous version
- **SSE transport binding**: wire format for streaming with reconnection support in [OAP-SSE](spec/oap-sse.md)

## Installation

```sh
go get github.com/urmzd/open-artifact-protocol
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

	oap "github.com/urmzd/open-artifact-protocol"
)

func main() {
	format := "text/html"
	body := `<oap:target id="msg">hello</oap:target>`
	raw, _ := json.Marshal(oap.SynthesizeContentItem{Body: body})

	envelope := oap.Envelope{
		Protocol: oap.ProtocolVersion,
		ID:       "greeting",
		Version:  1,
		Name:     oap.NameSynthesize,
		Meta:     oap.Meta{Format: &format},
		Content:  []json.RawMessage{raw},
	}

	artifact, handle, err := oap.Apply(nil, envelope)
	if err != nil {
		panic(err)
	}
	fmt.Println(artifact.Body)
	fmt.Println(handle.Name)
}
```

## Apply Engine

```go
artifact, handle, err := oap.Apply(baseArtifact, envelope)
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
store := oap.NewArtifactStore(10)
artifact, handle, err := store.Apply(envelope)
```

For edit envelopes, the store enforces `stored_version == envelope.version - 1`. A `synthesize` envelope can create or reset an artifact chain.

## Commands

| Recipe | Description |
|---|---|
| `just build` | Compile the Go package |
| `just test` | Run Go tests |
| `just evalset` | Regenerate SAIGE observations from the OAP eval corpus |
| `just report` | Regenerate eval economics and reliability report from `metrics.json` files |
| `just run` | Run live OpenAI-compatible eval experiments |
| `just check` | Run formatting, vet, and tests |

Equivalent direct commands:

```sh
go test ./...
go vet ./...
go run ./internal/evalsetgen
go run ./internal/evalreport
test -z "$(gofmt -l *.go cmd/oap-eval/*.go evalset/*.go internal/evalsetgen/*.go internal/evalreport/*.go internal/liveeval/*.go)"
```

## Evaluation Corpus

The eval corpus remains under [`assets/evals/`](assets/evals/). The runnable replacement is a Go live runner in [`cmd/oap-eval`](cmd/oap-eval/) plus the SAIGE-compatible eval set in [`assets/evals/saige/`](assets/evals/saige/), generated by [`evalset`](evalset/) as `github.com/urmzd/saige/eval` observations.

Read the benchmark numbers as a risk profile, not a blanket claim that OAP always wins. Raw `comparison` savings measure the protocol attempt itself. New live runs also report `reliability` and `economics` blocks: miss counts/rates, fallback-adjusted savings, retry tax, and init-inclusive amortized savings. If a OAP edit misses, production systems should expect to pay for the failed envelope attempt and then fall back to full regeneration; if an edit applies but corrupts content, only correctness oracles catch it.

The current `gpt-4o-mini` metrics set has 96 metrics files, of which 67 are comparable non-degenerate base-vs-OAP economics runs. In that comparable set, fallback-adjusted edit-token savings are 37.8%, init-inclusive fallback savings are 34.0%, and the miss rate is 14.9%. The strongest savings appear in artifacts above roughly 2 KB, repeated records/pages, dashboards, feeds, catalogs, API payloads, and code files with stable section boundaries. Tiny configs and one-off edits can lose because marker, supervisor inventory, and envelope overhead exceed the full-regeneration baseline.

The runner now treats target discovery as part of the supervisor loop. For non-JSON artifacts it inventories valid `<oap:target>` IDs with `list_targets()`. For JSON artifacts it inventories valid JSON Pointers with `list_paths()`. Generated envelopes are validated before apply; invented targets or invalid paths are rejected and repaired before they can become apply misses.

The previous Rust eval binary and benchmark harness have been removed in favor of the Go runner and SAIGE eval integration. Committed `metrics.json` and `results.md` files are measured artifacts: do not hand-edit measured numbers. Regenerate the SAIGE observation set with `just evalset` after changing experiment prompts or metadata, and regenerate measured result tables with `just report` after changing committed metrics.

## License

This project is dual-licensed:

- **Code** (Go files, build files): [Apache License 2.0](LICENSE)
- **Specification, eval assets, and documentation**: [CC-BY 4.0](LICENSE-CC-BY-4.0)

See [NOTICE](NOTICE) for details. Attribution is required under both licenses.
