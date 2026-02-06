package do

import (
	"os"
	"path/filepath"

	"github.com/brezzgg/go-packages/lg"
)

func (d *Do) GetDir(dir string) (string, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	stat, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", lg.Ef("not dir")
	}
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
