package servers

func GetServersFromPayload(opt map[string]interface{}) (servers []Server, err error) {
	var list []Server
	// TODO(ChrisMcKenzie): handle nil server list
	for _, val := range opt["list"].([]interface{}) {
		server := val.(map[interface{}]interface{})
		list = append(list, Server{
			Address: server["address"].(string),
			User:    server["username"].(string),
			// Password: server["password"].(string),
		})
	}

	servers = list

	return
}
