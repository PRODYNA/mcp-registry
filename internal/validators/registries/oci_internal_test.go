package registries

import "testing"

func TestIsAllowedRegistryAdditionalRegistries(t *testing.T) {
	tests := []struct {
		name                 string
		registry             string
		additionalRegistries []string
		want                 bool
	}{
		{
			name:                 "allows configured exact host",
			registry:             "registry.example.com",
			additionalRegistries: []string{"registry.example.com"},
			want:                 true,
		},
		{
			name:                 "ignores surrounding whitespace and host case",
			registry:             "registry.example.com",
			additionalRegistries: []string{" REGISTRY.EXAMPLE.COM "},
			want:                 true,
		},
		{
			name:                 "allows configured host and port",
			registry:             "registry.internal:5000",
			additionalRegistries: []string{"registry.internal:5000"},
			want:                 true,
		},
		{
			name:                 "does not allow custom subdomains",
			registry:             "sub.registry.example.com",
			additionalRegistries: []string{"registry.example.com"},
			want:                 false,
		},
		{
			name:                 "does not allow lookalike domains",
			registry:             "registry.example.com.attacker.test",
			additionalRegistries: []string{"registry.example.com"},
			want:                 false,
		},
		{
			name:     "built-in registries stay allowed",
			registry: "ghcr.io",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAllowedRegistry(tt.registry, tt.additionalRegistries); got != tt.want {
				t.Errorf("isAllowedRegistry(%q, %v) = %t, want %t", tt.registry, tt.additionalRegistries, got, tt.want)
			}
		})
	}
}
