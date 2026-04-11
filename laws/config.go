package laws

import (
	"github.com/franchb/fptest-go/engine"
	enginerapid "github.com/franchb/fptest-go/engine/rapid"
)

// Option configures law verification behavior.
type Option func(*config)

type config struct {
	runner engine.Runner
}

// WithRunner sets the PBT engine runner for law verification.
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
