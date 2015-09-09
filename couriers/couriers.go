package couriers

import "gopkg.in/yaml.v2"

type (
	Deployment struct {
		Command string `yaml:"command"`
		Servers struct {
			Provider string        `yaml:"provider"`
			Options  yaml.MapSlice `yaml:"options"`
		} `yaml:"servers"`
	}
)
