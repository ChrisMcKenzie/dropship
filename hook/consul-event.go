package hook

import (
	"encoding/json"
	"errors"
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

	name, ok := config["name"].(string)
	serv, ok := config["service"].(string)
	tag, ok := config["tag"].(string)
	node, ok := config["node"].(string)

	if !ok {
		return errors.New("Consul Hook: invalid config")
	}

	id, meta, err := client.Event().Fire(&api.UserEvent{
		Name:          name,
		Payload:       plBytes,
		ServiceFilter: serv,
		TagFilter:     tag,
		NodeFilter:    node,
	}, nil)

	fmt.Println(id, meta, err)

	return err
}
