package service

type Artifact struct {
	Type        string `hcl:",key"`
	Bucket      string `hcl:"bucket"`
	Path        string `hcl:"path"`
	Destination string `hcl:"destination"`
}

type Config struct {
	Name          string   `hcl:",key"`
	CheckInterval string   `hcl:"checkInterval"`
	PostCommand   string   `hcl:postCommand`
	PreCommand    string   `hcl:preCommand`
	Sequential    bool     `hcl:"sequentialUpdates"`
	Artifact      Artifact `hcl:"artifact,expand"`
}
