package i18n

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	bundleFiles []string
)

// loads translation file from yaml
func LoadFile(files ...string) error {
	bundle := make(Bundle)
	for i, f := range files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		var b Bundle
		err = yaml.Unmarshal(data, &b)
		if err != nil {
			return err
		}
		for key, messages := range b {
			if _, ok := bundle[key]; !ok {
				bundle[key] = make([]string, 0)
			}
			if bundleMessages, ok := bundle[key]; ok {
				if len(bundleMessages) < i {
					for j := 0; j < i; j++ {
						bundleMessages = append(bundleMessages, "")
					}
				}
				bundleMessages = append(bundleMessages, messages...)
				bundle[key] = bundleMessages
			}
		}
	}
	instance = translator{bundle: bundle}
	if devMode {
		bundleFiles = files
	}
	return nil
}
