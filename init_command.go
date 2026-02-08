package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "usage: specguard init")
		return 2
	}

	statuses, err := initSpecWorkspace("spec")
	if err != nil {
		var failure *initFailure
		if errors.As(err, &failure) {
			fmt.Printf("ERROR %s %s: %s\n", failure.target, failure.ruleKey, failure.message)
			return 2
		}

		fmt.Printf("ERROR spec execution-error: %v\n", err)
		return 2
	}

	for _, status := range statuses {
		if status.created {
			fmt.Printf("CREATED %s\n", status.path)
		} else {
			fmt.Printf("EXISTS %s\n", status.path)
		}
	}
	fmt.Println("OK spec init complete")
	return 0
}

func initSpecWorkspace(specRoot string) ([]initDirStatus, error) {
	requiredDirs := []string{
		filepath.Join(specRoot, "hls"),
		filepath.Join(specRoot, "lls"),
		filepath.Join(specRoot, "tc"),
		filepath.Join(specRoot, "shared"),
	}

	// Preflight to prevent partial writes when a required path is a file.
	for _, dir := range requiredDirs {
		info, err := os.Stat(dir)
		if err == nil {
			if !info.IsDir() {
				return nil, &initFailure{
					target:  filepath.ToSlash(dir),
					ruleKey: "path-conflict",
					message: "required path exists and is not a directory",
				}
			}
			continue
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		return nil, &initFailure{
			target:  filepath.ToSlash(dir),
			ruleKey: "stat-failed",
			message: err.Error(),
		}
	}

	statuses := make([]initDirStatus, 0, len(requiredDirs))
	for _, dir := range requiredDirs {
		info, err := os.Stat(dir)
		if err == nil {
			if !info.IsDir() {
				return nil, &initFailure{
					target:  filepath.ToSlash(dir),
					ruleKey: "path-conflict",
					message: "required path exists and is not a directory",
				}
			}
			statuses = append(statuses, initDirStatus{path: filepath.ToSlash(dir), created: false})
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, &initFailure{
				target:  filepath.ToSlash(dir),
				ruleKey: "stat-failed",
				message: err.Error(),
			}
		}

		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			return nil, &initFailure{
				target:  filepath.ToSlash(dir),
				ruleKey: "mkdir-failed",
				message: mkErr.Error(),
			}
		}
		statuses = append(statuses, initDirStatus{path: filepath.ToSlash(dir), created: true})
	}

	return statuses, nil
}
