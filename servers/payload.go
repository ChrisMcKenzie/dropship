package servers

import "encoding/json"

func GetServersFromPayload(opt json.RawMessage) (servers []Server, err error) {
	list := struct {
		List []Server `json:"list"`
	}{}
	if err = json.Unmarshal(opt, &list); err != nil {
		return
	}

	servers = list.List

	return
}
