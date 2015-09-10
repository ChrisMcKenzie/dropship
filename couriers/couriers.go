package couriers

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type (
	Deployment struct {
		Id       int      `yaml:"-"`
		Owner    string   `yaml:"-"`
		Repo     string   `yaml:"-"`
		Commands []string `yaml:"commands"`
		Servers  struct {
			Provider string                 `yaml:"provider"`
			Options  map[string]interface{} `yaml:"options"`
		} `yaml:"servers"`
	}

	Courier interface {
		Handle(*http.Request) (Deployment, error)
		UpdateStatus(Deployment, string, string) error
	}
)

func parseDeployment(file []byte) (Deployment, error) {
	log.Debug(string(file))
	var d Deployment
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}
