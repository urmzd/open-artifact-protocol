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

> [!WARNING]
> **This repository is archived and no longer maintained.** It is kept as a historical record of a design and of what was measured. Read [why it was archived](#status-archived) and the [evaluation limitations](#evaluation-limitations) before relying on it.

> [!NOTE]
> **Disclaimer.** This project is `v0`, and the protocol, schemas, and APIs were never stabilized. Nothing here is a claim about how any vendor's product works internally.

Open Artifact Protocol (OAP), formerly Generative Artifact Protocol (GAP), lets LLMs update structured text artifacts by emitting small JSON envelopes instead of regenerating full files. The native implementation is now Go-only: protocol types, `<oap:target>` marker resolution, JSON Pointer edits, a stateless apply engine, and a versioned in-memory store live in the root Go package.

## Status: archived

This repository was archived on 2026-10-10. The specification, the Go apply engine, and the eval corpus remain available under their existing licenses. The project is no longer developed.

The archive is not a statement that the idea failed: editing a stored artifact with small targeted operations is now common practice. It is not a claim to be an industry standard either: this repository records no independent implementation of the protocol, and the measurements here do not show that it beats current editing approaches. The project was archived because later tools solve most of the same practical problems, and the remaining open questions are research questions, not a reason to keep maintaining this code.

### Original motivation

The protocol was first published in February 2026 as Generative Artifact Protocol (GAP). It was designed around limits of the LLM systems available at the time:

- **Expensive regeneration.** Models rewrote a whole artifact to change one value, and output tokens were the main cost.
- **Context management.** Each regenerated version stayed in the conversation, so input grew with every edit.
- **Unreliable recall.** Reproducing existing text exactly, as a string-replace edit requires, often failed over long artifacts. The spec calls this out as "hallucinated search strings" in [Section 3.2](spec/oap.md#32-targets).
- **Structured artifacts across interactions.** There was no common way to keep a named, versioned artifact consistent across turns, agents, and tools.

[Section 7 of the spec](spec/oap.md#7-artifact-type-interface) states the same motivation as context bloat, token waste, and hallucination.

### Architectural intent

**Target IDs are addressable data points, not only a token optimization.** An `<oap:target>` ID names one value or region in the artifact: a metric, a price, a table row, a section. The IDs were meant as stable handles for retrieval, targeted updates, references, and validation. Cheaper edits were one use of them.

**A string-replace edit is a recall task. An ID edit is a lookup.** With string replacement, the model must reproduce the existing text exactly, and the match fails on any slip. With IDs, the model picks a short name from a list it was just shown.

**The orchestrator and maintain separation was intentional.** OAP separates an orchestrator that holds only a handle (ID, version, target list) from a short-lived maintain context that loads the artifact and edits it. The spec names this [context offloading](spec/oap.md#context-offloading). IDs are what make the separation possible, because an agent that never saw the body can still address a target by name. Quoting text requires seeing it.

The spec documents ID addressing and context offloading. The uses for references, validation, and goal-driven self-correction listed under [Remaining hypotheses](#remaining-hypotheses) are the author's design intent, recorded here on 2026-10-10. They were not specified or built.

### Industry convergence

| Date | Event |
|---|---|
| 2026-02-26 | First commit, as GAP |
| 2026-03-29 | `v0.1.0`, as GAP |
| 2026-04-02 | Spec draft 0.1 |
| 2026-06-18 | Claude Code Artifacts launches in beta |
| 2026-10-10 | Renamed to OAP (`v0.19.0`) and archived |

Claude Code Artifacts and other current tools now solve many of the same practical problems: the artifact body is stored instead of re-emitted, edits are small and targeted, and each publish is versioned. They do so with different implementation choices. As observed from the outside in one Claude Code session, a single agent writes the file, keeps it in its own context, and edits it by quoting existing text. OAP uses pre-placed IDs and a separate maintain context.

This is convergence on similar goals. It is not a claim of derivation, and the two designs are not architecturally equivalent. Models with long context have also largely outgrown the recall weakness, so string-replace editing now works reliably in practice.

### Remaining hypotheses

None of the following is established by this repository. They are possibilities that may be worth testing, particularly for smaller models and for complex stateful workflows.

- **Stable identifiers.** On small models, ID addressing may reach a higher exact-match rate than string-replace editing, with the difference closing as models get stronger.
- **Scoped context.** Editing from a handle and target list alone may keep input proportional to the number of targets instead of the size of the artifact.
- **References.** A classifier, a reviewer, or a second model may flag a problem by target ID, so the producer can act on exactly that region.
- **Deterministic validation.** A check may be written against a named data point (for example, "`total-revenue` equals the sum of the rows") instead of against a text pattern.
- **Declarative goals.** A goal may be stated as the expected state of named targets, so a producer can compare and correct target by target without re-judging the whole artifact.

### Evaluation limitations

- **The baseline is full regeneration.** The results under [`assets/evals/`](assets/evals/experiments/EXPERIMENT.md) compare OAP against regenerating the whole artifact. They do not compare OAP against string-replace editing, and they do not establish that OAP is better than current editing approaches.
- **The context-isolated workflow was not evaluated.** The live OAP flow sent the maintain context the full artifact on every edit, together with the target list. Editing from the handle and target list alone, which is the workflow the design was built for, was never run.
- **No comparison across model sizes.** The committed metrics come from one model.
- **One informal observation, not a benchmark.** In two edits watched in a Claude Code session on 2026-10-10, the quoted old text was roughly 2 to 3 percent of edit output. That suggests the output-token margin between ID edits and string-replace edits is small for small edits. It is two data points.

### Naming and the historical record

The project was renamed from Generative Artifact Protocol (GAP) to Open Artifact Protocol (OAP) on 2026-10-10, released as `v0.19.0`. The rename is breaking: the protocol string is now `oap/0.1`, markers are `<oap:target>`, the Go package is `oap`, and the module path is `github.com/urmzd/open-artifact-protocol`. Releases up to `v0.18.0` use the former `gap` identifiers and the former module path.

What was and was not renamed in the record:

- **`CHANGELOG.md`** keeps the original names.
- **Recorded model outputs** under `assets/evals/experiments/*/outputs/` and the fixtures under `assets/evals/apply-engine/` are kept exactly as the models produced them. They contain `<gap:target>` and `gap/0.1`, because that is what the models were asked for at the time.
- **Prompts, directory names, and metrics keys** in the eval corpus were renamed mechanically so the tooling still loads them. The measured numbers were not changed.
- **To reproduce a recorded run** exactly as it was executed, use the `v0.18.0` tag.

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

Read the benchmark numbers as a risk profile, not a blanket claim that OAP always wins. Raw `comparison` savings measure the protocol attempt itself. New live runs also report `reliability` and `economics` blocks: miss counts/rates, fallback-adjusted savings, retry tax, and init-inclusive amortized savings. If an OAP edit misses, production systems should expect to pay for the failed envelope attempt and then fall back to full regeneration; if an edit applies but corrupts content, only correctness oracles catch it.

The current `gpt-4o-mini` metrics set has 96 metrics files, of which 67 are comparable non-degenerate base-vs-OAP economics runs. In that comparable set, fallback-adjusted edit-token savings are 37.8%, init-inclusive fallback savings are 34.0%, and the miss rate is 14.9%. The strongest savings appear in artifacts above roughly 2 KB, repeated records/pages, dashboards, feeds, catalogs, API payloads, and code files with stable section boundaries. Tiny configs and one-off edits can lose because marker, supervisor inventory, and envelope overhead exceed the full-regeneration baseline.

The runner now treats target discovery as part of the supervisor loop. For non-JSON artifacts it inventories valid `<oap:target>` IDs with `list_targets()`. For JSON artifacts it inventories valid JSON Pointers with `list_paths()`. Generated envelopes are validated before apply; invented targets or invalid paths are rejected and repaired before they can become apply misses.

The previous Rust eval binary and benchmark harness have been removed in favor of the Go runner and SAIGE eval integration. Committed `metrics.json` and `results.md` files are measured artifacts: do not hand-edit measured numbers. Regenerate the SAIGE observation set with `just evalset` after changing experiment prompts or metadata, and regenerate measured result tables with `just report` after changing committed metrics.

## License

This project is dual-licensed:

- **Code** (Go files, build files): [Apache License 2.0](LICENSE)
- **Specification, eval assets, and documentation**: [CC-BY 4.0](LICENSE-CC-BY-4.0)

See [NOTICE](NOTICE) for details. Attribution is required under both licenses.
