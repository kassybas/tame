package loader

import (
	"github.com/kassybas/tame/internal/keywords"
	"github.com/kassybas/tame/internal/tamefile"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func readFile(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fc := string(b)
	return fc, nil
}


func Load(filePath string) (tamefile.Teafile, error) {
	fc, err := readFile(filePath)
	if err != nil {
		return tamefile.Teafile{}, err
	}

	t := tamefile.Teafile{}

	// Setting default opts here to be able to differentiate
	// between the empty opt (go default) and unset opt
	t.Sets.DefaultOptsContainer = keywords.OptsNotSet
	err = yaml.Unmarshal([]byte(fc), &t)
	return t, err
}
