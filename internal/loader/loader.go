package loader

import (
	"io/ioutil"

	"github.com/kassybas/mate/schema"
	"gopkg.in/yaml.v2"
)

func readFile(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fc := string(b)
	return fc, nil
}

func Load(filePath string) (schema.Tamefile, error) {
	fc, err := readFile(filePath)
	if err != nil {
		return schema.Tamefile{}, err
	}

	t := schema.Tamefile{}

	// Setting default opts here to be able to differentiate
	// between the empty opt (go default) and unset opt
	// t.Sets.DefaultOptsContainer = keywords.OptsNotSet
	err = yaml.UnmarshalStrict([]byte(fc), &t)
	return t, err
}
