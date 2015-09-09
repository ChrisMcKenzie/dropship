package servers

func GetServersFromPayload(opt map[string]interface{}) (servers []Server, err error) {
	var list []Server
	// TODO(ChrisMcKenzie): handle nil server list
	for _, val := range opt["list"].([]interface{}) {
		server := val.(map[string]string)
		list = append(list, Server{
			Address:  server["address"],
			User:     server["username"],
			Password: server["password"],
		})
	}

	servers = list

	return
}
