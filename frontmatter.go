package main

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

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
