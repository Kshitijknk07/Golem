package plugin

import (
	"Golem/internal/metrics"
	"context"
	"time"
)

type CheckPlugin interface {
	Name() string

	Type() metrics.HealthCheckType

	Description() string

	Execute(ctx context.Context, target string, timeout time.Duration) (metrics.HealthCheckStatus, string, time.Duration)

	ValidateConfig(config map[string]interface{}) error
}

type Registry struct {
	plugins map[string]CheckPlugin
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]CheckPlugin),
	}
}

func (r *Registry) Register(plugin CheckPlugin) {
	r.plugins[plugin.Name()] = plugin
}

func (r *Registry) Get(name string) (CheckPlugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

func (r *Registry) List() []CheckPlugin {
	result := make([]CheckPlugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		result = append(result, plugin)
	}
	return result
}
