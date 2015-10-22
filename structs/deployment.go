package structs

type Deployment struct {
	Services []Service `hcl:"service,expand"`
}

type Service struct {
	Id               string         `hcl:",key"`
	Name             string         `hcl:"name"`
	SequentialUpdate bool           `hcl:"sequentialUpdate"`
	CheckInterval    string         `hcl:"checkInterval"`
	Command          string         `hcl:"command"`
	Artifact         ArtifactConfig `hcl:"artifact,expand"`
	Hash             string         `hcl:"-"`
}

type ArtifactConfig struct {
	Repo   string `hcl:",key"`
	Bucket string `hcl:"bucket"`
	Type   string `hcl:"type"`
	Path   string `hcl:"path"`
	Dest   string `hcl:"destination"`
}
