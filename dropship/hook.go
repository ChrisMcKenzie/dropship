package dropship

// TemplateData specifies some system configurations that can be used in hook
// options using go template syntax
type TemplateData struct {
	Config
	Hostname string
}

// Hook is an interface that defines a method for executing custom
// behavior at certain hook points
type Hook interface {
	Execute(config HookConfig, service Config) error
}

// GetHookByName returns a hook by the id name given.
//
// TODO(ChrisMcKenzie): this is kind of ugly and could probably be cleaned
// through the use of hooks being packages and registering with a hook manager
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
