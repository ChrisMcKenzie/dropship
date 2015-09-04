package couriers

import "encoding/json"

type (
	Deployment struct {
		Command string `json:"command"`
		Servers struct {
			Provider string          `json:"provider"`
			Options  json.RawMessage `json:"options"`
		} `json:"servers"`
	}
)
