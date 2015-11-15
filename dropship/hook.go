package dropship

type TemplateData struct {
	Config
	Hostname string
}

type Hook interface {
	Execute(config HookConfig, service Config) error
}

func GetHookByName(name string) Hook {
	switch name {
	case "script":
		return ScriptHook{}
	case "consul-event":
		return ConsulEventHook{}
	case "graphite-event":
		return GraphiteEventHook{}
	}

	return nil
}
