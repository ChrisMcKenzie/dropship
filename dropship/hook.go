package dropship

var hooks = make(map[string]Hook)

func init() {
	registerHook(NewConsulEventHook())
	registerHook(NewGraphiteEventHook())
	registerHook(NewScriptHook())
}

// TemplateData specifies some system configurations that can be used in hook
// options using go template syntax
type TemplateData struct {
	Config
	Hostname string
}

type HookMeta struct {
	name string
}

// Name returns the name of this hook
func (h *HookMeta) Name() string {
	return h.name
}

// Hook is an interface that defines a method for executing custom
// behavior at certain hook points
type Hook interface {
	Execute(config HookConfig, service Config) error
	Name() string
}

func registerHook(hook Hook) {
	hooks[hook.Name()] = hook
}

// GetHookByName returns a hook by the id name given.
//
// TODO(ChrisMcKenzie): this is kind of ugly and could probably be cleaned
// through the use of hooks being packages and registering with a hook manager
func GetHookByName(name string) Hook {
	if hook, ok := hooks[name]; ok {
		return hook
	}
	return nil
}
