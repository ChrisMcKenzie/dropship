package servers

func GetServersFromPayload(opt map[string]interface{}) (servers []Server, err error) {
	var list []Server
	// TODO(ChrisMcKenzie): handle nil server list
	for _, val := range opt["list"].([]string) {
		list = append(list, Server{Address: val})
	}

	servers = list

	return
}
