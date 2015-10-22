package agent

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ChrisMcKenzie/dropship/structs"
	"github.com/hashicorp/hcl"
)

func loadServices(root string) (d []structs.Service) {
	files, _ := filepath.Glob(root + "/*.hcl")
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		var deploy structs.Deployment
		hcl.Decode(&deploy, string(data))
		d = append(d, deploy.Services...)
	}
	return
}
