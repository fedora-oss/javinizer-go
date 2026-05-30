package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/scraperutil"
	"gopkg.in/yaml.v3"
)

// Test to see what unmarshal receives
func TestDebugUnmarshalReceives(t *testing.T) {
	yamlData := `
scrapers:
    user_agent: test
    r18dev:
        enabled: true
`
	type Wrapper struct {
		Scrapers *ScrapersConfig `yaml:"scrapers"`
	}

	wrapper := &Wrapper{
		Scrapers: &ScrapersConfig{},
	}

	err := yaml.Unmarshal([]byte(yamlData), wrapper)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	t.Logf("After unmarshal:")
	t.Logf("  wrapper.Scrapers.Overrides: %d keys", len(wrapper.Scrapers.Overrides))
	t.Logf("  wrapper.Scrapers.UserAgent: %q", wrapper.Scrapers.UserAgent)
}

// Test the full unmarshal path with actual config
func TestDebugFullUnmarshalPath(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "configs", "config.yaml")
	data, err := os.ReadFile(repoRoot)
	if err != nil {
		t.Skipf("configs/config.yaml not found, skipping integration test")
	}

	// First, unmarshal to find the scrapers section
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal doc error: %v", err)
	}

	// Find scrapers node
	var scrapersNode *yaml.Node
	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		key := doc.Content[0].Content[i].Value
		if key == "scrapers" {
			scrapersNode = doc.Content[0].Content[i+1]
			break
		}
	}

	if scrapersNode == nil {
		t.Fatal("scrapers node not found")
	}

	t.Logf("Scrapers node kind: %d", scrapersNode.Kind)

	// Check what keys are in scrapers
	keyCount := 0
	for i := 0; i < len(scrapersNode.Content); i += 2 {
		keyCount++
		if keyCount <= 15 {
			t.Logf("  scrapers key[%d]: %s", i, scrapersNode.Content[i].Value)
		}
	}
	t.Logf("  ... (total %d keys)", keyCount)

	// Now manually unmarshal the scrapers section
	scrapersData, err := yaml.Marshal(scrapersNode)
	if err != nil {
		t.Fatalf("yaml.Marshal scrapersNode error: %v", err)
	}
	t.Logf("\nScrapers YAML:\n%s", string(scrapersData))

	// Test what the factory returns for known scrapers
	t.Logf("\nChecking factories:")
	for _, name := range []string{"r18dev", "dmm", "javlibrary", "javdb"} {
		factory := scraperutil.GetConfigFactory(name)
		if factory == nil {
			t.Logf("  %s: factory is NIL", name)
		} else {
			t.Logf("  %s: factory found", name)
		}

		flattenFn := scraperutil.GetFlattenFunc(name)
		if flattenFn == nil {
			t.Logf("  %s: flattenFunc is NIL", name)
		} else {
			t.Logf("  %s: flattenFunc found", name)
		}
	}

	// Now unmarshal into ScrapersConfig directly
	cfg := &Config{
		Scrapers: ScrapersConfig{},
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		t.Fatalf("yaml.Unmarshal cfg error: %v", err)
	}

	t.Logf("\nAfter Load:")
	t.Logf("  UserAgent: %q", cfg.Scrapers.UserAgent)
	t.Logf("  Proxy.Enabled: %v", cfg.Scrapers.Proxy.Enabled)
	t.Logf("  Proxy.Profiles: %+v", cfg.Scrapers.Proxy.Profiles)
	t.Logf("  Overrides: %d keys", len(cfg.Scrapers.Overrides))
}

// Test the ScrapersConfig.UnmarshalYAML with detailed logging
func TestDebugScrapersUnmarshalYAML(t *testing.T) {
	yamlData := `
scrapers:
    user_agent: test-agent
    priority:
        - r18dev
        - dmm
    proxy:
        enabled: true
        default_profile: main
        profiles:
            main:
                url: "http://proxy.example.com:8080"
    r18dev:
        enabled: true
        rate_limit: 500
    javlibrary:
        enabled: true
        language: en
`

	sc := &ScrapersConfig{}

	// Manually call UnmarshalYAML by using yaml.Unmarshal with a wrapper
	type Wrapper struct {
		Scrapers *ScrapersConfig `yaml:"scrapers"`
	}
	wrapper := &Wrapper{Scrapers: sc}

	err := yaml.Unmarshal([]byte(yamlData), wrapper)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	t.Logf("After unmarshal:")
	t.Logf("  UserAgent: %q", sc.UserAgent)
	t.Logf("  Priority: %v", sc.Priority)
	t.Logf("  Proxy.Enabled: %v", sc.Proxy.Enabled)
	t.Logf("  Proxy.Profiles: %+v", sc.Proxy.Profiles)
	t.Logf("  Overrides: %d keys", len(sc.Overrides))
	for k, v := range sc.Overrides {
		t.Logf("    %s: %+v", k, v)
	}
}

// Test FlatToScraperConfig directly
func TestDebugFlatToScraperConfig(t *testing.T) {
	// Get factory for r18dev
	factory := scraperutil.GetConfigFactory("r18dev")
	if factory == nil {
		t.Skip("r18dev factory is nil (scraper packages not imported in this test package)")
	}

	// Create concrete type
	concrete := factory()

	// Marshal sample YAML
	sampleYAML := `
enabled: true
rate_limit: 500
user_agent: test
`
	err := yaml.Unmarshal([]byte(sampleYAML), concrete)
	if err != nil {
		t.Fatalf("Unmarshal concrete error: %v", err)
	}

	// Try to flatten
	result := FlatToScraperConfig("r18dev", concrete)
	if result == nil {
		t.Logf("FlatToScraperConfig returned nil")
	} else {
		t.Logf("FlatToScraperConfig result: %+v", result)
	}
}

// Simulate what UnmarshalYAML does
func TestDebugSimulateUnmarshalYAML(t *testing.T) {
	yamlData := `
scrapers:
    user_agent: test-agent
    r18dev:
        enabled: true
        rate_limit: 500
    javlibrary:
        enabled: true
        language: en
`

	var generic map[string]any
	err := yaml.Unmarshal([]byte(yamlData), &generic)
	if err != nil {
		t.Fatalf("Unmarshal generic error: %v", err)
	}

	t.Logf("generic keys: %v", fmt.Sprintf("%v", generic))

	scrapers := generic["scrapers"]
	if scrapers == nil {
		t.Fatal("scrapers key not found in generic")
	}

	scrapersMap, ok := scrapers.(map[string]any)
	if !ok {
		t.Fatalf("scrapers is not map[string]any, it's %T", scrapers)
	}

	t.Logf("scrapersMap keys: %v", fmt.Sprintf("%v", scrapersMap))

	// Now try to marshal and unmarshal one scraper
	r18devValue := scrapersMap["r18dev"]
	data, err := yaml.Marshal(r18devValue)
	if err != nil {
		t.Fatalf("yaml.Marshal r18devValue error: %v", err)
	}
	t.Logf("\nr18dev YAML:\n%s", string(data))

	// Get factory
	factory := scraperutil.GetConfigFactory("r18dev")
	if factory == nil {
		t.Skip("r18dev factory is nil (scraper packages not imported in this test package)")
	}

	concrete := factory()

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	err = decoder.Decode(concrete)
	if err != nil {
		t.Fatalf("Decode concrete error: %v", err)
	}

	t.Logf("concrete after decode: %+v", concrete)

	// Flatten
	result := FlatToScraperConfig("r18dev", concrete)
	if result == nil {
		t.Fatal("FlatToScraperConfig returned nil")
	}
	t.Logf("result: %+v", result)
}
