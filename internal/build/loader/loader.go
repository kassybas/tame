package loader

import (
	"io/ioutil"

	"github.com/kassybas/tame/schema"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func ReadFile(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Load(filePath string) (schema.Tamefile, map[string]interface{}, error) {
	fc, err := ReadFile(filePath)
	if err != nil {
		return schema.Tamefile{}, nil, err
	}

	var raw map[string]interface{}
	err = yaml.UnmarshalStrict(fc, &raw)
	if err != nil {
		return schema.Tamefile{}, nil, err
	}
	var md mapstructure.Metadata
	var result schema.Tamefile
	err = mapstructure.DecodeMetadata(raw, &result, &md)
	dynamic := make(map[string]interface{})
	// collect dynamic key-values
	for _, k := range md.Unused {
		dynamic[k] = raw[k]
	}
	return result, dynamic, err
}
