package service

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

type Artifact struct {
	Type        string `hcl:",key"`
	Bucket      string `hcl:"bucket"`
	Path        string `hcl:"path"`
	Destination string `hcl:"destination"`
}

type Hook map[string]map[string]interface{}

type Config struct {
	Name          string     `hcl:",key"`
	CheckInterval string     `hcl:"checkInterval"`
	PostCommand   string     `hcl:"postCommand"`
	PreCommand    string     `hcl:"preCommand"`
	BeforeHooks   []Hook     `hcl:"before"`
	AfterHooks    []Hook     `hcl:"after"`
	Sequential    bool       `hcl:"sequentialUpdates"`
	Artifact      []Artifact `hcl:"artifact"`
	Hash          string
}

type ServiceFile struct {
	Services []Config `hcl:"service"`
}

func LoadServices(root string) (d []Config, err error) {
	files, _ := filepath.Glob(root + "/*.hcl")
	for _, file := range files {
		data, err := readFile(file)
		if err != nil {
			return nil, err
		}

		var deploy ServiceFile
		err = hcl.Decode(&deploy, data)
		if err != nil {
			return nil, err
		}

		d = append(d, deploy.Services...)
	}
	return
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
