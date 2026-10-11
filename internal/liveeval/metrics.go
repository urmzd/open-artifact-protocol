package liveeval

import (
	"time"

	"github.com/urmzd/saige/eval/harness"
)

// Harness types reused where the JSON wire format is identical to OAP's
// original local definitions. The OAP-specific per-turn keys
// (envelope_parsed, apply_succeeded, envelope_name) travel in
// TurnResult.Extra and are flattened into the same JSON object.
type (
	TurnMetrics = harness.TurnMetrics
	TurnResult  = harness.TurnResult
	FlowMetrics = harness.FlowMetrics
	Comparison  = harness.Comparison
	Reliability = harness.Reliability
)

type Metrics struct {
	ExperimentID  string          `json:"experiment_id"`
	Model         string          `json:"model"`
	Provider      string          `json:"provider"`
	Timestamp     string          `json:"timestamp"`
	Format        string          `json:"format"`
	BaseTurn0     *TurnMetrics    `json:"base_turn0,omitempty"`
	OAPTurn0      *TurnMetrics    `json:"oap_turn0,omitempty"`
	StatelessT0   *TurnMetrics    `json:"stateless_turn0,omitempty"`
	StatelessFlow *FlowMetrics    `json:"stateless_flow,omitempty"`
	DefaultFlow   *FlowMetrics    `json:"default_flow,omitempty"`
	OAPFlow       *OapFlowMetrics `json:"oap_flow,omitempty"`
	Comparison    *Comparison     `json:"comparison,omitempty"`
	Decomposition *Decomposition  `json:"decomposition,omitempty"`
	Reliability   *Reliability    `json:"reliability,omitempty"`
	Economics     *Economics      `json:"economics,omitempty"`
	Validity      *Validity       `json:"validity,omitempty"`
}

type OapFlowMetrics struct {
	FlowMetrics
	EnvelopeParseRate float64 `json:"envelope_parse_rate"`
	ApplySuccessRate  float64 `json:"apply_success_rate"`
}

type Decomposition struct {
	InputSavingsBVsAPct  float64 `json:"input_savings_b_vs_a_pct"`
	OutputSavingsCVsBPct float64 `json:"output_savings_c_vs_b_pct"`
	InputSavingsCVsAPct  float64 `json:"input_savings_c_vs_a_pct"`
	OutputSavingsCVsAPct float64 `json:"output_savings_c_vs_a_pct"`
}

type Economics struct {
	FallbackAssumption           string              `json:"fallback_assumption"`
	MeasuredTotalTokenSavingsPct float64             `json:"measured_total_token_savings_pct"`
	FallbackAdjusted             *FallbackAdjusted   `json:"fallback_adjusted,omitempty"`
	Amortized                    *AmortizedEconomics `json:"amortized,omitempty"`
}

type FallbackAdjusted struct {
	InputTokens               uint64  `json:"input_tokens"`
	OutputTokens              uint64  `json:"output_tokens"`
	TotalTokens               uint64  `json:"total_tokens"`
	MissAttemptInputTokens    uint64  `json:"miss_attempt_input_tokens"`
	MissAttemptOutputTokens   uint64  `json:"miss_attempt_output_tokens"`
	MissAttemptTotalTokens    uint64  `json:"miss_attempt_total_tokens"`
	FallbackRetryInputTokens  uint64  `json:"fallback_retry_input_tokens"`
	FallbackRetryOutputTokens uint64  `json:"fallback_retry_output_tokens"`
	FallbackRetryTotalTokens  uint64  `json:"fallback_retry_total_tokens"`
	InputTokenSavingsPct      float64 `json:"input_token_savings_pct"`
	OutputTokenSavingsPct     float64 `json:"output_token_savings_pct"`
	TotalTokenSavingsPct      float64 `json:"total_token_savings_pct"`
}

type AmortizedEconomics struct {
	EditTurns                            int     `json:"edit_turns"`
	BaseInitInclusiveTokens              uint64  `json:"base_init_inclusive_tokens"`
	OAPInitInclusiveTokens               uint64  `json:"oap_init_inclusive_tokens"`
	FallbackInitInclusiveTokens          uint64  `json:"fallback_init_inclusive_tokens"`
	MeasuredInitInclusiveTokenSavingsPct float64 `json:"measured_init_inclusive_token_savings_pct"`
	FallbackInitInclusiveTokenSavingsPct float64 `json:"fallback_init_inclusive_token_savings_pct"`
}

type Validity struct {
	OAPRunDegenerate  bool `json:"oap_run_degenerate"`
	BaseInputMonotone bool `json:"base_input_monotone"`
}

// assembleMetrics maps the per-flow harness results into the OAP Metrics
// document, then fills the derived analytics.
func assembleMetrics(model string, exp harness.Script, results map[string]harness.FlowResult) Metrics {
	metrics := Metrics{
		ExperimentID: exp.ID,
		Model:        model,
		Provider:     "openai-compatible",
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Format:       exp.Format,
	}
	if result, ok := results[baseFlowName]; ok {
		t0 := result.Turn0
		flow := harness.ToFlowMetrics(result.Turns)
		metrics.BaseTurn0 = &t0
		metrics.DefaultFlow = &flow
	}
	if result, ok := results[statelessFlowName]; ok {
		t0 := result.Turn0
		flow := harness.ToFlowMetrics(result.Turns)
		metrics.StatelessT0 = &t0
		metrics.StatelessFlow = &flow
	}
	if result, ok := results[oapFlowName]; ok {
		t0 := result.Turn0
		oapMetrics := OapFlowMetrics{
			FlowMetrics: harness.ToFlowMetrics(result.Turns),
			EnvelopeParseRate: harness.Rate(result.Turns, func(t TurnResult) bool {
				return extraBool(t.Extra, "envelope_parsed")
			}),
			ApplySuccessRate: harness.Rate(result.Turns, func(t TurnResult) bool {
				return extraBool(t.Extra, "apply_succeeded")
			}),
		}
		metrics.OAPTurn0 = &t0
		metrics.OAPFlow = &oapMetrics
	}
	fillDerived(&metrics)
	return metrics
}

func fillDerived(metrics *Metrics) {
	if metrics.DefaultFlow != nil && metrics.OAPFlow != nil {
		comparison := harness.Compare(*metrics.DefaultFlow, metrics.OAPFlow.FlowMetrics)
		metrics.Comparison = &comparison
		metrics.Economics = economics(metrics.BaseTurn0, metrics.OAPTurn0, metrics.DefaultFlow, metrics.OAPFlow)
	}
	if metrics.OAPFlow != nil {
		metrics.Reliability = harness.ComputeReliability(metrics.OAPFlow.PerTurn)
	}
	if metrics.DefaultFlow != nil && metrics.StatelessFlow != nil && metrics.OAPFlow != nil {
		aIn, aOut := initInclusive(metrics.BaseTurn0, metrics.DefaultFlow)
		bIn, bOut := initInclusive(metrics.StatelessT0, metrics.StatelessFlow)
		cIn, cOut := initInclusive(metrics.OAPTurn0, &metrics.OAPFlow.FlowMetrics)
		metrics.Decomposition = &Decomposition{
			InputSavingsBVsAPct:  harness.Pct(aIn, bIn),
			OutputSavingsCVsBPct: harness.Pct(bOut, cOut),
			InputSavingsCVsAPct:  harness.Pct(aIn, cIn),
			OutputSavingsCVsAPct: harness.Pct(aOut, cOut),
		}
	}
	if metrics.DefaultFlow != nil || metrics.OAPFlow != nil {
		metrics.Validity = &Validity{
			OAPRunDegenerate: oapDegenerate(metrics.OAPFlow),
			BaseInputMonotone: metrics.DefaultFlow == nil ||
				inputMonotone(metrics.DefaultFlow.PerTurn),
		}
	}
}

func economics(baseTurn0 *TurnMetrics, oapTurn0 *TurnMetrics, baseFlow *FlowMetrics, oapFlow *OapFlowMetrics) *Economics {
	baseIn, baseOut, oapIn, oapOut := alignedTotals(baseFlow.PerTurn, oapFlow.PerTurn)
	fallbackIn, fallbackOut := oapIn, oapOut
	var missAttemptIn, missAttemptOut, retryIn, retryOut uint64
	for i, oapTurn := range oapFlow.PerTurn {
		if i >= len(baseFlow.PerTurn) || !oapTurn.Failed {
			continue
		}
		baseTurn := baseFlow.PerTurn[i]
		missAttemptIn += oapTurn.InputTokens
		missAttemptOut += oapTurn.OutputTokens
		retryIn += baseTurn.InputTokens
		retryOut += baseTurn.OutputTokens
		fallbackIn += baseTurn.InputTokens
		fallbackOut += baseTurn.OutputTokens
	}

	baseTotal := baseIn + baseOut
	oapTotal := oapIn + oapOut
	fallbackTotal := fallbackIn + fallbackOut
	econ := &Economics{
		FallbackAssumption:           "on each missed OAP edit, run the baseline full-regeneration edit after the failed OAP attempt",
		MeasuredTotalTokenSavingsPct: harness.Pct(baseTotal, oapTotal),
		FallbackAdjusted: &FallbackAdjusted{
			InputTokens:               fallbackIn,
			OutputTokens:              fallbackOut,
			TotalTokens:               fallbackTotal,
			MissAttemptInputTokens:    missAttemptIn,
			MissAttemptOutputTokens:   missAttemptOut,
			MissAttemptTotalTokens:    missAttemptIn + missAttemptOut,
			FallbackRetryInputTokens:  retryIn,
			FallbackRetryOutputTokens: retryOut,
			FallbackRetryTotalTokens:  retryIn + retryOut,
			InputTokenSavingsPct:      harness.Pct(baseIn, fallbackIn),
			OutputTokenSavingsPct:     harness.Pct(baseOut, fallbackOut),
			TotalTokenSavingsPct:      harness.Pct(baseTotal, fallbackTotal),
		},
	}
	if baseTurn0 != nil && oapTurn0 != nil {
		baseSession := baseTotal + baseTurn0.InputTokens + baseTurn0.OutputTokens
		oapSession := oapTotal + oapTurn0.InputTokens + oapTurn0.OutputTokens
		fallbackSession := fallbackTotal + oapTurn0.InputTokens + oapTurn0.OutputTokens
		econ.Amortized = &AmortizedEconomics{
			EditTurns:                            len(oapFlow.PerTurn),
			BaseInitInclusiveTokens:              baseSession,
			OAPInitInclusiveTokens:               oapSession,
			FallbackInitInclusiveTokens:          fallbackSession,
			MeasuredInitInclusiveTokenSavingsPct: harness.Pct(baseSession, oapSession),
			FallbackInitInclusiveTokenSavingsPct: harness.Pct(baseSession, fallbackSession),
		}
	}
	return econ
}

func alignedTotals(baseTurns []TurnResult, oapTurns []TurnResult) (uint64, uint64, uint64, uint64) {
	var baseIn, baseOut, oapIn, oapOut uint64
	limit := len(baseTurns)
	if len(oapTurns) < limit {
		limit = len(oapTurns)
	}
	for i := 0; i < limit; i++ {
		baseIn += baseTurns[i].InputTokens
		baseOut += baseTurns[i].OutputTokens
		oapIn += oapTurns[i].InputTokens
		oapOut += oapTurns[i].OutputTokens
	}
	return baseIn, baseOut, oapIn, oapOut
}

func initInclusive(t0 *TurnMetrics, flow *FlowMetrics) (uint64, uint64) {
	in, out := flow.TotalInputTokens, flow.TotalOutputTokens
	if t0 != nil {
		in += t0.InputTokens
		out += t0.OutputTokens
	}
	return in, out
}

func oapDegenerate(flow *OapFlowMetrics) bool {
	if flow == nil || len(flow.PerTurn) < 2 {
		return false
	}
	first := flow.PerTurn[0].OutputBytes
	for _, turn := range flow.PerTurn[1:] {
		if turn.OutputBytes != first {
			return false
		}
	}
	return true
}

func inputMonotone(turns []TurnResult) bool {
	for i := 1; i < len(turns); i++ {
		if turns[i].InputTokens < turns[i-1].InputTokens {
			return false
		}
	}
	return true
}

// extraBool reads a bool key from a TurnResult Extra map, false when absent.
func extraBool(extra map[string]any, key string) bool {
	value, ok := extra[key].(bool)
	return ok && value
}
