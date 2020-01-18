package lint

import yaml "gopkg.in/yaml.v3"

func marshalhack(src, dst interface{}) error {
	b, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, dst)
}
