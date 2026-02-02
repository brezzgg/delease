package parser

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/brezzgg/go-packages/lg"
)

var confRe = regexp.MustCompile(`delease\.yaml$`)

func FindConfig(def, cwd string) ([]byte, error) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return nil, lg.Ef("get cwd: %w", err)
		}
	}

	if def != "" {
		isAbs := filepath.IsAbs(def)
		if !isAbs {
			def = filepath.Join(cwd, def)
		}
		b, err := os.ReadFile(def)
		if err != nil {
			return nil, lg.Ef("failed to read file '%s': %w", def, err)
		}
		return b, nil
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil, lg.Ef("failed to read cwd: %w", err)
	}

	fileName := ""
	for _, enrty := range entries {
		if enrty.IsDir() {
			continue
		}
		if confRe.MatchString(enrty.Name()) {
			fileName = filepath.Join(cwd, enrty.Name())
			break
		}
	}

	if fileName == "" {
		return nil, lg.Ef("config not found, specify it use -C, or change working directory use -D")
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, lg.Ef("failed to read file '%s': %w", fileName, err)
	}

	return b, nil
}
