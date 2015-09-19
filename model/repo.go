package model

type Repo struct {
	Model
	User     User   `json:"user"`
	UserId   int64  `json:"-"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Courier  string `json:"courier"`
	URL      string `json:"url"`
	CloneURL string `json:"clone_url"`
	Active   bool   `json:"active"`
	Private  bool   `json:"-"`
	Token    string `json:"-"`
}
