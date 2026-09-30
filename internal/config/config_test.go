package config

import (
	"reflect"
	"testing"
)

func TestNewConfigAdditionalOCIRegistries(t *testing.T) {
	t.Setenv("MCP_REGISTRY_ADDITIONAL_OCI_REGISTRIES", " Registry.Example.com, registry.internal:5000,,registry.example.com")

	cfg := NewConfig()

	want := []string{"registry.example.com", "registry.internal:5000"}
	if !reflect.DeepEqual(cfg.AdditionalOCIRegistries, want) {
		t.Fatalf("AdditionalOCIRegistries = %v, want %v", cfg.AdditionalOCIRegistries, want)
	}
}

func TestNewConfigAdditionalOCIRegistriesEmpty(t *testing.T) {
	t.Setenv("MCP_REGISTRY_ADDITIONAL_OCI_REGISTRIES", "")

	if got := NewConfig().AdditionalOCIRegistries; len(got) != 0 {
		t.Fatalf("AdditionalOCIRegistries = %v, want empty", got)
	}
}

func TestNormalizeOCIRegistryHostsRejectsInvalid(t *testing.T) {
	for _, host := range []string{"https://registry.example.com", "registry.example.com/repo", "*.example.com", "user@registry.example.com", "a b"} {
		if _, err := normalizeOCIRegistryHosts([]string{host}); err == nil {
			t.Errorf("normalizeOCIRegistryHosts(%q) succeeded, want error", host)
		}
	}
}
