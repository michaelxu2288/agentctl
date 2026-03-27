package tools

import "fmt"

type Tool interface {
	Name() string
	Description() string
	Call(params map[string]any) (string, error)
}

type Registry struct {
	items map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{items: map[string]Tool{}}
}

func (r *Registry) Register(t Tool) {
	r.items[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, error) {
	t, ok := r.items[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return t, nil
}

func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.items))
	for n := range r.items {
		out = append(out, n)
	}
	return out
}
