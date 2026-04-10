package prop

import (
	"github.com/franchb/fptest/engine"
	enginerapid "github.com/franchb/fptest/engine/rapid"
)

// Option configures property test behavior.
type Option func(*config)

type config struct {
	runner engine.Runner
}

// WithRunner sets the PBT engine runner for property tests.
func WithRunner(r engine.Runner) Option {
	return func(c *config) { c.runner = r }
}

func resolveConfig(opts []Option) config {
	cfg := config{runner: enginerapid.RapidRunner{}}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
