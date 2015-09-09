package couriers

import "gopkg.in/yaml.v2"

type (
	Deployment struct {
		Command string `yaml:"command"`
		Servers struct {
			Provider string                 `yaml:"provider"`
			Options  map[string]interface{} `yaml:"options"`
		} `yaml:"servers"`
	}
)

func parseDeployment(file []byte) (Deployment, error) {
	log.Debug(file)
	var d Deployment
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}
