package loader

import (
	"fmt"
	"io/ioutil"

	"github.com/kassybas/mate/internal/keywords"
	"github.com/kassybas/mate/scheme"
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

func Load(filePath string) (scheme.Tamefile, error) {
	fmt.Println("Loading", filePath)
	fc, err := readFile(filePath)
	if err != nil {
		return scheme.Tamefile{}, err
	}

	t := scheme.Tamefile{}

	// Setting default opts here to be able to differentiate
	// between the empty opt (go default) and unset opt
	t.Sets.DefaultOptsContainer = keywords.OptsNotSet
	err = yaml.Unmarshal([]byte(fc), &t)
	fmt.Printf("%+v\n", t.Targets)
	panic("ok")
	return t, err
}
