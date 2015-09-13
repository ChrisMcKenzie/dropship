package servers

import "github.com/hashicorp/consul/api"

func GetServersFromConsul(options map[string]interface{}) (servers []Server, err error) {

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return
	}

	catalog := client.Catalog()
	services, _, err := catalog.Service(options["service_name"].(string), "", nil)
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
