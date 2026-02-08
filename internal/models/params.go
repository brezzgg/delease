package models

import (
	"regexp"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

type ParamsParser struct {
	paramRe *regexp.Regexp
	taskRe  *regexp.Regexp
}

func NewParamsParser() *ParamsParser {
	return &ParamsParser{
		paramRe: regexp.MustCompile(`([A-Za-z0-9_-]+)=(?:"([^"]*)"|([\sA-Za-z0-9_-]+))`),
		taskRe:  regexp.MustCompile(`^(?P<task>[A-Za-z0-9_-]+):(?P<params>.*)$`),
	}
}

func (p *ParamsParser) Parse(input string) (string, *VarSource, error) {
	if !strings.Contains(input, ":") {
		return input, nil, nil
	}

	matches := p.taskRe.FindStringSubmatch(input)

	if matches == nil {
		return "", nil, lg.Ef("invalid format")
	}

	taskName := matches[1]
	paramsStr := matches[2]

	if taskName == "" {
		return "", nil, lg.Ef("task name cannot be empty")
	}

	paramMatches := p.paramRe.FindAllStringSubmatch(paramsStr, -1)

	if len(paramMatches) == 0 && strings.TrimSpace(paramsStr) != "" {
		return "", nil, lg.Ef("invalid parameters format")
	}

	data := make(map[string]*Var, len(paramMatches))

	for _, match := range paramMatches {
		varName := match[1]
		varValue := ""

		if match[2] != "" {
			varValue = match[2]
		} else {
			varValue = match[3]
		}

		if varName == "" {
			return "", nil, lg.Ef("var name cannot be empty")
		}

		if varValue == "" {
			return "", nil, lg.Ef("var value cannot be empty for '%s'", varName)
		}

		if _, exists := data[varName]; exists {
			return "", nil, lg.Ef("duplicate variable name: %s", varName)
		}

		data[varName] = NewVar(varValue, VarTypeDynamic)
	}

	src := &VarSource{}
	src.SetSource(data)

	return taskName, src, nil
}
