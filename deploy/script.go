package deploy

type (
	Deploy struct {
		Servers map[string]Server `yaml:"servers"`
		Tasks   map[string][]Step `yaml:"tasks"`
	}

	Server struct {
		Provider string                 `yaml:"provider"`
		Options  map[string]interface{} `yaml:"options"`
	}

	Step        map[string]StepOptions
	StepOptions map[string]string
)
