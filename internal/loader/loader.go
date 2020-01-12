package loader

import (
	"io/ioutil"

	"github.com/kassybas/tame/schema"
	"gopkg.in/yaml.v2"
)

func readFile(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Load(filePath string) (schema.Tamefile, error) {
	fc, err := readFile(filePath)
	if err != nil {
		return schema.Tamefile{}, err
	}

	t := schema.Tamefile{}

	err = yaml.UnmarshalStrict(fc, &t)
	return t, err
}
