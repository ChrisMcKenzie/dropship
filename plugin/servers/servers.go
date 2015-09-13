package servers

type (
	Server struct {
		Address  string `json:"address"`
		User     string `json:"username"`
		Password string `json:"password"`
	}
)
