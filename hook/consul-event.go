package hook

import (
	"encoding/json"
	"fmt"

	"github.com/ChrisMcKenzie/dropship/service"
	"github.com/hashicorp/consul/api"
)

type ConsulEventHook struct{}

func (h ConsulEventHook) Execute(config map[string]interface{}, service service.Config) error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	payload := map[string]string{
		"hash": service.Hash,
	}

	plBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	id, meta, err := client.Event().Fire(&api.UserEvent{
		Name:          config["name"].(string),
		Payload:       plBytes,
		ServiceFilter: config["service"].(string),
		TagFilter:     config["tag"].(string),
		NodeFilter:    config["node"].(string),
	}, nil)

	fmt.Println(id, meta, err)

	return err
}
