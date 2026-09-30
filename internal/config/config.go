package config

import (
	"fmt"
	"strings"

	env "github.com/caarlos0/env/v11"
)

// Config holds the application configuration
// See .env.example for more documentation
type Config struct {
	ServerAddress            string   `env:"SERVER_ADDRESS" envDefault:":8080"`
	DatabaseURL              string   `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/mcp-registry?sslmode=disable"`
	SeedFrom                 string   `env:"SEED_FROM" envDefault:""`
	Version                  string   `env:"VERSION" envDefault:"dev"`
	GithubClientID           string   `env:"GITHUB_CLIENT_ID" envDefault:""`
	GithubClientSecret       string   `env:"GITHUB_CLIENT_SECRET" envDefault:""`
	JWTPrivateKey            string   `env:"JWT_PRIVATE_KEY" envDefault:""`
	EnableAnonymousAuth      bool     `env:"ENABLE_ANONYMOUS_AUTH" envDefault:"false"`
	EnableRegistryValidation bool     `env:"ENABLE_REGISTRY_VALIDATION" envDefault:"true"`
	AdditionalOCIRegistries  []string `env:"ADDITIONAL_OCI_REGISTRIES" envSeparator:","`

	GitHubOIDCAudience string `env:"GITHUB_OIDC_AUDIENCE" envDefault:""`

	// OIDC Configuration
	OIDCEnabled      bool   `env:"OIDC_ENABLED" envDefault:"false"`
	OIDCIssuer       string `env:"OIDC_ISSUER" envDefault:""`
	OIDCClientID     string `env:"OIDC_CLIENT_ID" envDefault:""`
	OIDCExtraClaims  string `env:"OIDC_EXTRA_CLAIMS" envDefault:""`
	OIDCEditPerms    string `env:"OIDC_EDIT_PERMISSIONS" envDefault:""`
	OIDCPublishPerms string `env:"OIDC_PUBLISH_PERMISSIONS" envDefault:""`
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
		Prefix: "MCP_REGISTRY_",
	})
	if err != nil {
		panic(err)
	}
	cfg.AdditionalOCIRegistries, err = normalizeOCIRegistryHosts(cfg.AdditionalOCIRegistries)
	if err != nil {
		panic(err)
	}
	return &cfg
}

// normalizeOCIRegistryHosts trims, lowercases and de-duplicates registry hosts, dropping
// empty entries. It rejects anything that is not a bare host (optionally with a port),
// such as URLs or wildcard patterns, because those would silently never match.
func normalizeOCIRegistryHosts(hosts []string) ([]string, error) {
	var normalized []string
	seen := make(map[string]bool)
	for _, host := range hosts {
		host = strings.ToLower(strings.TrimSpace(host))
		if host == "" {
			continue
		}
		if strings.ContainsAny(host, "/*@?# \t") {
			return nil, fmt.Errorf("invalid MCP_REGISTRY_ADDITIONAL_OCI_REGISTRIES entry %q: must be a hostname with optional port (no scheme, path or wildcard)", host)
		}
		if !seen[host] {
			seen[host] = true
			normalized = append(normalized, host)
		}
	}
	return normalized, nil
}
