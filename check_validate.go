package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func validateSpecWorkspace(specRoot string) ([]violation, error) {
	requiredDirs := []specDir{
		{path: filepath.Join(specRoot, "hls"), prefix: "HLS-"},
		{path: filepath.Join(specRoot, "lls"), prefix: "LLS-"},
		{path: filepath.Join(specRoot, "tc"), prefix: "TC-"},
	}

	var out []violation
	var files []specFile
	seenIDs := map[string]string{}

	for _, dir := range requiredDirs {
		info, err := os.Stat(dir.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				out = append(out, violation{
					target:  filepath.ToSlash(dir.path),
					ruleKey: "missing-directory",
					message: "required directory does not exist",
				})
				continue
			}
			return nil, fmt.Errorf("stat %s: %w", dir.path, err)
		}
		if !info.IsDir() {
			out = append(out, violation{
				target:  filepath.ToSlash(dir.path),
				ruleKey: "missing-directory",
				message: "required directory path is not a directory",
			})
			continue
		}

		err = filepath.WalkDir(dir.path, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(d.Name()) != ".md" {
				return nil
			}

			filePath := filepath.ToSlash(path)
			fileID, linkIDs, fileViolations, err := validateSpecFile(path, dir.prefix)
			if err != nil {
				return err
			}
			out = append(out, fileViolations...)

			if fileID != "" {
				if prevPath, exists := seenIDs[fileID]; exists {
					out = append(out, violation{
						target:  filePath,
						ruleKey: "duplicate-id",
						message: fmt.Sprintf("id %s already declared in %s", fileID, prevPath),
					})
				} else {
					seenIDs[fileID] = filePath
				}

				files = append(files, specFile{
					path:    filePath,
					id:      fileID,
					linkIDs: linkIDs,
				})
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk %s: %w", dir.path, err)
		}
	}

	for _, f := range files {
		for _, linkedID := range f.linkIDs {
			if !idPattern.MatchString(linkedID) {
				continue
			}
			if _, exists := seenIDs[linkedID]; !exists {
				out = append(out, violation{
					target:  f.path,
					ruleKey: "unresolved-link",
					message: fmt.Sprintf("links contains unknown id %s", linkedID),
				})
			}
		}
	}

	return out, nil
}

func validateSpecFile(path, expectedPrefix string) (string, []string, []violation, error) {
	var out []violation
	filePath := filepath.ToSlash(path)
	base := filepath.Base(path)
	re := fileNamePatternByPrefix[expectedPrefix]
	matches := re.FindStringSubmatch(base)
	expectedID := ""
	if matches == nil {
		out = append(out, violation{
			target:  filePath,
			ruleKey: "invalid-filename",
			message: fmt.Sprintf("expected filename pattern %s###-<slug>.md", strings.TrimSuffix(expectedPrefix, "-")),
		})
	} else {
		expectedID = matches[1]
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return "", nil, nil, fmt.Errorf("read %s: %w", path, err)
	}
	content := string(contentBytes)
	frontMatter, body, err := extractFrontMatter(content)
	if err != nil {
		out = append(out, violation{
			target:  filePath,
			ruleKey: "front-matter-missing",
			message: err.Error(),
		})
		return "", nil, out, nil
	}

	id, links, parseViolations := validateFrontMatter(frontMatter, expectedID)
	for i := range parseViolations {
		parseViolations[i].target = filePath
	}
	out = append(out, parseViolations...)

	if expectedPrefix == "HLS-" && !containsGherkinFence(body) {
		out = append(out, violation{
			target:  filePath,
			ruleKey: "missing-gherkin",
			message: "expected at least one fenced gherkin code block",
		})
	}

	return id, links, out, nil
}
