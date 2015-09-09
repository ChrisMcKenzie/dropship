package couriers

import "gopkg.in/yaml.v2"

type (
	Deployment struct {
		Commands []string `yaml:"commands"`
		Servers  struct {
			Provider string                 `yaml:"provider"`
			Options  map[string]interface{} `yaml:"options"`
		} `yaml:"servers"`
	}
)

func parseDeployment(file []byte) (Deployment, error) {
	var d Deployment
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}
