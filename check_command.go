package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "usage: specguard check")
		return 2
	}

	violations, err := validateSpecWorkspace("spec")
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %v\n", err)
		return 2
	}

	if len(violations) == 0 {
		fmt.Println("OK spec check passed")
		return 0
	}

	sort.Slice(violations, func(i, j int) bool {
		if violations[i].target != violations[j].target {
			return violations[i].target < violations[j].target
		}
		if violations[i].ruleKey != violations[j].ruleKey {
			return violations[i].ruleKey < violations[j].ruleKey
		}
		return violations[i].message < violations[j].message
	})

	for _, v := range violations {
		fmt.Printf("ERROR %s %s: %s\n", v.target, v.ruleKey, v.message)
	}

	return 1
}
