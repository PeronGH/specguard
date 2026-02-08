package main

import "regexp"

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

type initDirStatus struct {
	path    string
	created bool
}

type initFailure struct {
	target  string
	ruleKey string
	message string
}

func (e *initFailure) Error() string {
	return e.message
}
