package compliance

import (
	"fmt"
	"strings"

	"github.com/scim2/test-suite/scim"
	"github.com/scim2/test-suite/spec"
)

func featureSupported(cfg Config, f spec.Feature) bool {
	if cfg.Force[f] {
		return true
	}
	switch f {
	case spec.Core:
		return true
	case spec.Filter:
		return cfg.Features.Filter
	case spec.Patch:
		return cfg.Features.Patch
	case spec.Bulk:
		return cfg.Features.Bulk
	case spec.Sort:
		return cfg.Features.Sort
	case spec.ETag:
		return cfg.Features.ETag
	case spec.ChangePassword:
		return cfg.Features.ChangePassword
	case spec.Users:
		return cfg.Features.Users
	case spec.Groups:
		return cfg.Features.Groups
	case spec.TestResource:
		return cfg.Features.TestResource
	case spec.Me:
		return false
	}
	return false
}

func printResult(name string, outcome Outcome, msg string, logs []string) {
	symbol := "PASS"
	switch outcome {
	case Fail:
		symbol = "FAIL"
	case Warn:
		symbol = "WARN"
	case Skip:
		symbol = "SKIP"
	}
	fmt.Printf("  [%s] %s\n", symbol, name)
	if msg != "" {
		fmt.Printf("         %s\n", msg)
	}
	for _, l := range logs {
		fmt.Printf("         %s\n", l)
	}
}

// Config holds the settings for a compliance run.
type Config struct {
	Client   *scim.Client
	Features *scim.Features
	Force    map[spec.Feature]bool
	Filter   string
	Verbose  bool
}

// Outcome represents the result of a compliance check.
type Outcome int

const (
	Pass Outcome = iota
	Fail
	Warn
	Skip
)

// Report holds the collected results of a compliance run.
type Report struct {
	Results []Result
}

// RunAll executes all tests attached to requirements and returns the results.
func RunAll(cfg Config) Report {
	var results []Result

	for i := range spec.All {
		req := &spec.All[i]
		if !req.Testable || len(req.Tests) == 0 {
			continue
		}

		if cfg.Filter != "" && !strings.HasPrefix(req.ID, cfg.Filter) {
			continue
		}

		adv := featureSupported(cfg, req.Feature)

		for _, test := range req.Tests {
			r := &spec.Run{
				Client:   cfg.Client,
				Features: cfg.Features,
			}
			r.Execute(test.Fn)

			// If the test used subtests, expand each into its own result.
			if subs := r.SubResults(); len(subs) > 0 {
				for _, sub := range subs {
					subName := test.Name + "/" + sub.Name
					name := req.ID + "/" + subName

					var outcome Outcome
					var msg string
					switch {
					case sub.Skipped:
						outcome = Skip
					case sub.Failed && adv:
						outcome = Fail
						msg = strings.Join(sub.Msgs, "; ")
					case sub.Failed && !adv:
						outcome = Warn
						msg = strings.Join(sub.Msgs, "; ")
					default:
						outcome = Pass
					}

					results = append(results, Result{
						RequirementIDs: []string{req.ID},
						TestName:       subName,
						Outcome:        outcome,
						Message:        msg,
					})

					if cfg.Verbose {
						printResult(name, outcome, msg, sub.Logs)
					}
				}
				continue
			}

			// No subtests: single result as before.
			name := req.ID + "/" + test.Name
			var outcome Outcome
			var msg string

			switch {
			case r.Skipped():
				outcome = Skip
			case r.Failed() && adv:
				outcome = Fail
				msg = strings.Join(r.Messages(), "; ")
			case r.Failed() && !adv:
				outcome = Warn
				msg = strings.Join(r.Messages(), "; ")
			default:
				outcome = Pass
			}

			results = append(results, Result{
				RequirementIDs: []string{req.ID},
				TestName:       test.Name,
				Outcome:        outcome,
				Message:        msg,
			})

			if cfg.Verbose {
				printResult(name, outcome, msg, r.Logs())
			}
		}
	}

	return Report{Results: results}
}

// Result links a test outcome to one or more requirement IDs.
type Result struct {
	RequirementIDs []string
	TestName       string
	Outcome        Outcome
	Message        string
}
