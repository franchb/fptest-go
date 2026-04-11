package gen

import (
	"github.com/franchb/fptest-go/engine"
	enginerapid "github.com/franchb/fptest-go/engine/rapid"
)

// ToEngine converts a Gen[A] to an engine.Generator[A] via the rapid adapter.
func ToEngine[A any](g Gen[A]) engine.Generator[A] {
	return enginerapid.Wrap(ToRapid(g))
}
