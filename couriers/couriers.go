package couriers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	Deployment struct {
		Id          int      `yaml:"-"`
		Owner       string   `yaml:"-"`
		Repo        string   `yaml:"-"`
		Environment string   `yaml:"-"`
		Commands    []string `yaml:"commands"`
		Servers     map[string]struct {
			Provider string                 `yaml:"provider"`
			Options  map[string]interface{} `yaml:"options"`
		} `yaml:"servers"`
	}

	Courier interface {
		Handle(*http.Request) (*Deployment, error)
		UpdateStatus(Deployment, string, string) error
	}
)

func ParseDeployment(file []byte) (Deployment, error) {
	log.Debug(string(file))
	var d Deployment
	err := yaml.Unmarshal(file, &d)
	if err != nil {
		return d, err
	}

	return d, nil
}
