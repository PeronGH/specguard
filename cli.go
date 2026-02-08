package main

import (
	"fmt"
	"os"
)

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: specguard <check|init>")
		return 2
	}

	switch args[0] {
	case "check":
		return runCheck(args[1:])
	case "init":
		return runInit(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "usage: specguard <check|init>")
		return 2
	}
}
