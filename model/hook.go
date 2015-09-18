package model

type Hook struct {
	Owner       string
	Repo        string
	Sha         string
	Branch      string
	PullRequest string
	Author      string
	Gravatar    string
	Timestamp   string
	Message     string
}
