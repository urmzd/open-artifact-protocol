// Package liveeval runs OAP's live eval experiments on top of the generic
// saige eval harness. The harness owns the chat client, the base and
// stateless flows, corpus skipping, and output plumbing; this package keeps
// only the protocol-specific parts: the OAP flow (synthesis plus envelope
// edit turns with repair loops), the supervisor target inventory, and the
// OAP metrics document schema.
package liveeval

import (
	"context"
	"fmt"

	"github.com/urmzd/open-artifact-protocol/evalset"
	"github.com/urmzd/saige/eval/harness"
)

// Flow names as spelled in metrics assembly and output directories.
const (
	baseFlowName      = "base"
	statelessFlowName = "stateless"
	oapFlowName       = "oap"
)

// System prompt keys used when mapping evalset experiments onto
// harness.Script.Systems.
const (
	baseSystemKey        = baseFlowName
	oapInitSystemKey     = "oap.init"
	oapMaintainSystemKey = "oap.maintain"
)

type Config struct {
	ExperimentsDir string
	Count          int
	IDFilter       string
	Flow           string
	Model          string
	APIBase        string
	APIKey         string
	Force          bool
}

func Run(ctx context.Context, cfg Config) error {
	flows, err := flowsFor(cfg.Flow)
	if err != nil {
		return err
	}
	inputs, err := evalset.LoadExperiments(cfg.ExperimentsDir)
	if err != nil {
		return err
	}
	experiments := harness.FilterScripts(toHarnessExperiments(inputs), cfg.IDFilter, cfg.Count)
	// NewClient omits temperature for models whose catalog row rejects it
	// at the default reasoning effort, including the gpt-6 family.
	client := harness.NewClient(cfg.APIBase, cfg.APIKey, cfg.Model)
	runner, err := harness.New(harness.Config{
		Client: client,
		Flows:  flows,
		Force:  cfg.Force,
		Assemble: func(exp harness.Script, results map[string]harness.FlowResult) (any, error) {
			return assembleMetrics(client.Model, exp, results), nil
		},
	})
	if err != nil {
		return err
	}
	return runner.Run(ctx, experiments)
}

func flowsFor(flow string) ([]harness.Flow, error) {
	switch flow {
	case baseFlowName:
		return []harness.Flow{harness.BaseFlow{}}, nil
	case statelessFlowName:
		return []harness.Flow{harness.StatelessFlow{}}, nil
	case oapFlowName:
		return []harness.Flow{oapFlow{}}, nil
	case "both":
		return []harness.Flow{harness.BaseFlow{}, oapFlow{}}, nil
	case "abc", "all":
		return []harness.Flow{harness.BaseFlow{}, harness.StatelessFlow{}, oapFlow{}}, nil
	default:
		return nil, fmt.Errorf("unknown flow %q", flow)
	}
}

func toHarnessExperiments(inputs []evalset.ExperimentInput) []harness.Script {
	experiments := make([]harness.Script, 0, len(inputs))
	for _, input := range inputs {
		turns := make([]harness.Turn, 0, len(input.Turns))
		for _, turn := range input.Turns {
			turns = append(turns, harness.Turn{Index: turn.Turn, Prompt: turn.Prompt})
		}
		experiments = append(experiments, harness.Script{
			ID:     input.ExperimentID,
			Format: input.Format,
			Dir:    input.Paths.ExperimentDir,
			Systems: map[string]string{
				baseSystemKey:        input.BaseSystem,
				oapInitSystemKey:     input.OAPInitSystem,
				oapMaintainSystemKey: input.OAPMaintain,
			},
			Turns: turns,
		})
	}
	return experiments
}
