package loader

import (
	"io/ioutil"

	"github.com/kassybas/tame/schema"
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

	err = yaml.UnmarshalStrict([]byte(fc), &t)
	return t, err
}
