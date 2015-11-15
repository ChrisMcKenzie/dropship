package hook

import "github.com/ChrisMcKenzie/dropship/service"

type TemplateData struct {
	service.Config
	Hostname string
}

type Hook interface {
	Execute(config service.HookConfig, service service.Config) error
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
