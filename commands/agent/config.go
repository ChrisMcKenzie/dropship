package agent

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type Config struct {
	ConfigDir string `hcl:"config"`
	LockHost  string `hcl:"lockhost"`
	Rackspace struct {
		User   string `hcl:"user"`
		Key    string `hcl:"key"`
		Region string `hcl:"region"`
	} `hcl:"rackspace"`
}

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = hcl.Decode(config, string(data))

	return config, err
}
