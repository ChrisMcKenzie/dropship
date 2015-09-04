package servers

import (
	"encoding/json"

	"github.com/hashicorp/consul/api"
)

func GetServersFromConsul(opt json.RawMessage) (servers []Server, err error) {
	options := struct {
		serviceName string `json:"service_name"`
	}{}
	err = json.Unmarshal(opt, &options)
	if err != nil {
		return
	}

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return
	}

	catalog := client.Catalog()
	services, _, err := catalog.Service(options.serviceName, "", nil)
	if err != nil {
		return
	}

	for _, service := range services {
		servers = append(
			servers,
			Server{
				Address: service.Address,
			},
		)
	}

	return
}
