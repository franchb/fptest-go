package hegelgen

import (
	"time"

	"github.com/franchb/fptest/engine"
	fpthegel "github.com/franchb/fptest/hegel"
	hegellib "hegel.dev/go/hegel"
)

// Emails returns a generator that produces email addresses.
func Emails() engine.Generator[string] { return fpthegel.Wrap(hegellib.Emails()) }

// URLs returns a generator that produces URLs.
func URLs() engine.Generator[string] { return fpthegel.Wrap(hegellib.URLs()) }

// Dates returns a generator that produces date values.
func Dates() engine.Generator[time.Time] { return fpthegel.Wrap(hegellib.Dates()) }

// Datetimes returns a generator that produces datetime values.
func Datetimes() engine.Generator[time.Time] { return fpthegel.Wrap(hegellib.Datetimes()) }

// Text returns a generator that produces string values of length in [minSize, maxSize].
func Text(minSize, maxSize int) engine.Generator[string] {
	return fpthegel.Wrap(hegellib.Text(minSize, maxSize))
}

// Booleans returns a generator that produces boolean values.
func Booleans() engine.Generator[bool] { return fpthegel.Wrap(hegellib.Booleans()) }

// FromRegex returns a generator that produces strings matching the given regex pattern.
func FromRegex(pattern string, fullmatch bool) engine.Generator[string] {
	return fpthegel.Wrap(hegellib.FromRegex(pattern, fullmatch))
}

// Integers returns a generator that produces integer values in [minVal, maxVal].
func Integers[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}](minVal, maxVal T) engine.Generator[T] {
	return fpthegel.Wrap(hegellib.Integers(minVal, maxVal))
}
