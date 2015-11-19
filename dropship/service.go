package dropship

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

type Artifact map[string]string
type HookConfig map[string]string
type HookDefinition map[string]HookConfig

type Config struct {
	Name          string              `hcl:",key"`
	CheckInterval string              `hcl:"checkInterval"`
	PostCommand   string              `hcl:"postCommand"`
	PreCommand    string              `hcl:"preCommand"`
	BeforeHooks   []HookDefinition    `hcl:"before"`
	AfterHooks    []HookDefinition    `hcl:"after"`
	Sequential    bool                `hcl:"sequentialUpdates"`
	RawArtifact   map[string]Artifact `hcl:"artifact,expand"`
	Artifact      Artifact            `hcl:"-"`
	Hash          string              `hcl:"hash"`
	Updater       Updater             `hcl:"-"`
	Locker        Locker              `hcl:"-"`
}

type ServiceFile struct {
	Services []Config `hcl:"service,expand"`
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

		for i, service := range deploy.Services {
			for key, value := range service.RawArtifact {
				value["type"] = key
				deploy.Services[i].Artifact = value
			}
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
