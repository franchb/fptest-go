package examples_test

import (
	"testing"

	"github.com/franchb/fptest-go/laws"
	"pgregory.net/rapid"
)

// --- Domain types ---

type TLSConfig struct {
	CertPath string
	KeyPath  string
	Enabled  bool
}

type ServerConfig struct {
	Host string
	Port int
	TLS  TLSConfig
}

type DBConfig struct {
	Host     string
	Port     int
	MaxConns int
}

type FeatureFlags struct {
	DarkMode   bool
	BetaAccess bool
}

type AppConfig struct {
	Server   ServerConfig
	Database DBConfig
	Features FeatureFlags
}

// --- Generators ---

func genTLSConfig() *rapid.Generator[TLSConfig] {
	return rapid.Custom(func(t *rapid.T) TLSConfig {
		return TLSConfig{
			CertPath: rapid.StringMatching(`/etc/certs/[a-z]+\.pem`).Draw(t, "cert"),
			KeyPath:  rapid.StringMatching(`/etc/certs/[a-z]+\.key`).Draw(t, "key"),
			Enabled:  rapid.Bool().Draw(t, "enabled"),
		}
	})
}

func genServerConfig() *rapid.Generator[ServerConfig] {
	return rapid.Custom(func(t *rapid.T) ServerConfig {
		return ServerConfig{
			Host: rapid.StringMatching(`[a-z]+\.example\.com`).Draw(t, "host"),
			Port: rapid.IntRange(1024, 65535).Draw(t, "port"),
			TLS:  genTLSConfig().Draw(t, "tls"),
		}
	})
}

func genDBConfig() *rapid.Generator[DBConfig] {
	return rapid.Custom(func(t *rapid.T) DBConfig {
		return DBConfig{
			Host:     rapid.StringMatching(`db-[a-z]+\.internal`).Draw(t, "host"),
			Port:     rapid.IntRange(1024, 65535).Draw(t, "port"),
			MaxConns: rapid.IntRange(1, 100).Draw(t, "maxConns"),
		}
	})
}

func genFeatureFlags() *rapid.Generator[FeatureFlags] {
	return rapid.Custom(func(t *rapid.T) FeatureFlags {
		return FeatureFlags{
			DarkMode:   rapid.Bool().Draw(t, "darkMode"),
			BetaAccess: rapid.Bool().Draw(t, "betaAccess"),
		}
	})
}

func genAppConfig() *rapid.Generator[AppConfig] {
	return rapid.Custom(func(t *rapid.T) AppConfig {
		return AppConfig{
			Server:   genServerConfig().Draw(t, "server"),
			Database: genDBConfig().Draw(t, "db"),
			Features: genFeatureFlags().Draw(t, "features"),
		}
	})
}

// --- Tests ---

// TestServerHostLens verifies the lens laws for AppConfig.Server.Host.
// A lawful lens guarantees that getting then setting is a no-op, you get back
// what you set, and setting twice is the same as setting once. Without these
// properties, configuration update helpers can silently lose data.
func TestServerHostLens(t *testing.T) {
	laws.LensLaws(t,
		genAppConfig(),
		rapid.StringMatching(`[a-z]+\.example\.com`),
		func(a, b AppConfig) bool { return a == b },
		func(a, b string) bool { return a == b },
		func(c AppConfig) string { return c.Server.Host },
		func(host string) func(AppConfig) AppConfig {
			return func(c AppConfig) AppConfig {
				c.Server.Host = host
				return c
			}
		},
	)
}

// TestTLSEnabledLens verifies lens laws for a 3-level nested field:
// AppConfig.Server.TLS.Enabled. Deeply nested setters are particularly prone to
// bugs where an intermediate struct is copied incorrectly.
func TestTLSEnabledLens(t *testing.T) {
	laws.LensLaws(t,
		genAppConfig(),
		rapid.Bool(),
		func(a, b AppConfig) bool { return a == b },
		func(a, b bool) bool { return a == b },
		func(c AppConfig) bool { return c.Server.TLS.Enabled },
		func(enabled bool) func(AppConfig) AppConfig {
			return func(c AppConfig) AppConfig {
				c.Server.TLS.Enabled = enabled
				return c
			}
		},
	)
}

// TestDBMaxConnsLens verifies lens laws for AppConfig.Database.MaxConns.
// The generator constrains MaxConns to [1, 100], matching a real-world domain
// constraint. Lens laws hold regardless of the value range.
func TestDBMaxConnsLens(t *testing.T) {
	laws.LensLaws(t,
		genAppConfig(),
		rapid.IntRange(1, 100),
		func(a, b AppConfig) bool { return a == b },
		func(a, b int) bool { return a == b },
		func(c AppConfig) int { return c.Database.MaxConns },
		func(maxConns int) func(AppConfig) AppConfig {
			return func(c AppConfig) AppConfig {
				c.Database.MaxConns = maxConns
				return c
			}
		},
	)
}

// TestFeatureFlagLens verifies lens laws for AppConfig.Features.DarkMode.
// Feature flags are commonly toggled via configuration updates — lens laws
// ensure the toggle operation is well-behaved.
func TestFeatureFlagLens(t *testing.T) {
	laws.LensLaws(t,
		genAppConfig(),
		rapid.Bool(),
		func(a, b AppConfig) bool { return a == b },
		func(a, b bool) bool { return a == b },
		func(c AppConfig) bool { return c.Features.DarkMode },
		func(darkMode bool) func(AppConfig) AppConfig {
			return func(c AppConfig) AppConfig {
				c.Features.DarkMode = darkMode
				return c
			}
		},
	)
}
