package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	fileNamePatternByPrefix = map[string]*regexp.Regexp{
		"HLS-": regexp.MustCompile(`^(HLS-\d{3})-[a-z0-9][a-z0-9-]*\.md$`),
		"LLS-": regexp.MustCompile(`^(LLS-\d{3})-[a-z0-9][a-z0-9-]*\.md$`),
		"TC-":  regexp.MustCompile(`^(TC-\d{3})-[a-z0-9][a-z0-9-]*\.md$`),
	}
	idPattern = regexp.MustCompile(`^(HLS|LLS|TC)-\d{3}$`)
)

type specDir struct {
	path   string
	prefix string
}

type violation struct {
	target  string
	ruleKey string
	message string
}

type specFile struct {
	path    string
	id      string
	linkIDs []string
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: specguard check")
		return 2
	}

	switch args[0] {
	case "check":
		return runCheck(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "usage: specguard check")
		return 2
	}
}

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

func extractFrontMatter(content string) (string, string, error) {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", "", errors.New("expected YAML front matter starting with --- on first line")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return "", "", errors.New("expected YAML front matter closing delimiter ---")
	}

	yml := strings.Join(lines[1:end], "\n")
	body := strings.Join(lines[end+1:], "\n")
	return yml, body, nil
}

func validateFrontMatter(frontMatter, expectedID string) (string, []string, []violation) {
	var out []violation

	var raw map[string]any
	if err := yaml.Unmarshal([]byte(frontMatter), &raw); err != nil {
		out = append(out, violation{
			ruleKey: "front-matter-invalid",
			message: fmt.Sprintf("invalid YAML front matter: %v", err),
		})
		return "", nil, out
	}

	idValue, idExists := raw["id"]
	if !idExists {
		out = append(out, violation{
			ruleKey: "id-missing",
			message: "front matter field id is required",
		})
	}
	id, ok := idValue.(string)
	if idExists && (!ok || strings.TrimSpace(id) == "") {
		out = append(out, violation{
			ruleKey: "id-invalid",
			message: "front matter field id must be a non-empty string",
		})
	}
	if ok && !idPattern.MatchString(id) {
		out = append(out, violation{
			ruleKey: "id-invalid",
			message: fmt.Sprintf("front matter id %s must match (HLS|LLS|TC)-###", id),
		})
	}
	if expectedID != "" && ok && id != expectedID {
		out = append(out, violation{
			ruleKey: "id-mismatch",
			message: fmt.Sprintf("front matter id %s does not match filename %s", id, expectedID),
		})
	}

	titleValue, titleExists := raw["title"]
	if !titleExists {
		out = append(out, violation{
			ruleKey: "title-missing",
			message: "front matter field title is required",
		})
	} else {
		title, ok := titleValue.(string)
		if !ok || strings.TrimSpace(title) == "" {
			out = append(out, violation{
				ruleKey: "title-invalid",
				message: "front matter field title must be a non-empty string",
			})
		}
	}

	statusValue, statusExists := raw["status"]
	if !statusExists {
		out = append(out, violation{
			ruleKey: "status-missing",
			message: "front matter field status is required",
		})
	} else {
		status, ok := statusValue.(string)
		if !ok {
			out = append(out, violation{
				ruleKey: "status-invalid",
				message: "front matter field status must be a string",
			})
		} else if status != "draft" && status != "active" && status != "deprecated" {
			out = append(out, violation{
				ruleKey: "status-invalid",
				message: fmt.Sprintf("front matter status %s is invalid", status),
			})
		}
	}

	var linkIDs []string
	if linksValue, hasLinks := raw["links"]; hasLinks {
		linksMap, ok := linksValue.(map[string]any)
		if !ok {
			out = append(out, violation{
				ruleKey: "links-invalid",
				message: "front matter field links must be an object of ID lists",
			})
		} else {
			for key, linked := range linksMap {
				items, ok := linked.([]any)
				if !ok {
					out = append(out, violation{
						ruleKey: "links-invalid",
						message: fmt.Sprintf("links.%s must be a list", key),
					})
					continue
				}
				for _, item := range items {
					id, ok := item.(string)
					if !ok || strings.TrimSpace(id) == "" {
						out = append(out, violation{
							ruleKey: "links-invalid",
							message: fmt.Sprintf("links.%s must contain non-empty string IDs", key),
						})
						continue
					}
					if !idPattern.MatchString(id) {
						out = append(out, violation{
							ruleKey: "links-invalid",
							message: fmt.Sprintf("links.%s contains invalid id %s", key, id),
						})
						continue
					}
					linkIDs = append(linkIDs, id)
				}
			}
		}
	}

	return id, linkIDs, out
}

func containsGherkinFence(content string) bool {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "```gherkin" {
			continue
		}
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "```" {
				return true
			}
		}
	}
	return false
}
