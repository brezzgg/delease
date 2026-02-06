package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

var confRe = regexp.MustCompile(`delease\.yaml$`)

func FindConfig(def, cwd string) ([]byte, string, error) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return nil, "", lg.Ef("get cwd: %w", err)
		}
	}

	if def != "" {
		isAbs := filepath.IsAbs(def)
		if !isAbs {
			def = filepath.Join(cwd, def)
		}
		b, err := os.ReadFile(def)
		if err != nil {
			return nil, "", lg.Ef("failed to read file '%s': %w", def, err)
		}
		return b, def, nil
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil, "", lg.Ef("failed to read cwd: %w", err)
	}

	var files []string
	for _, enrty := range entries {
		if enrty.IsDir() {
			continue
		}
		if confRe.MatchString(enrty.Name()) {
			files = append(files, filepath.Join(cwd, enrty.Name()))
		}
	}

	fileName := ""
	for _, v := range files {
		if strings.HasSuffix(v, "delease.yaml") {
			fileName = v
		}
	}
	if fileName == "" && len(files) > 0 {
		fileName = files[0]
	}

	if fileName == "" {
		return nil, "", lg.Ef("config not found, specify it use -C, or change working directory use -D")
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, "", lg.Ef("failed to read file '%s': %w", fileName, err)
	}

	return b, fileName, nil
}
