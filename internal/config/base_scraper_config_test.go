package config

import (
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/scraperutil"
	"github.com/stretchr/testify/assert"
)

func TestBaseScraperConfig_Interface(t *testing.T) {
	var _ scraperutil.ScraperConfigInterface = (*BaseScraperConfig)(nil)
}

func TestBaseScraperConfig_IsEnabled(t *testing.T) {
	testCases := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{"enabled true", true, true},
		{"enabled false", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{Enabled: tc.enabled}
			assert.Equal(t, tc.expected, c.IsEnabled())
		})
	}
}

func TestBaseScraperConfig_GetUserAgent(t *testing.T) {
	testCases := []struct {
		name      string
		userAgent string
		expected  string
	}{
		{"set user agent", "test-ua", "test-ua"},
		{"empty user agent", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{UserAgent: tc.userAgent}
			assert.Equal(t, tc.expected, c.GetUserAgent())
		})
	}
}

func TestBaseScraperConfig_GetRequestDelay(t *testing.T) {
	testCases := []struct {
		name         string
		requestDelay int
		expected     int
	}{
		{"delay 500", 500, 500},
		{"delay 0", 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{RequestDelay: tc.requestDelay}
			assert.Equal(t, tc.expected, c.GetRequestDelay())
		})
	}
}

func TestBaseScraperConfig_GetMaxRetries(t *testing.T) {
	testCases := []struct {
		name       string
		maxRetries int
		expected   int
	}{
		{"retries 3", 3, 3},
		{"retries 0", 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{MaxRetries: tc.maxRetries}
			assert.Equal(t, tc.expected, c.GetMaxRetries())
		})
	}
}

func TestBaseScraperConfig_GetProxy(t *testing.T) {
	proxy := &ProxyConfig{Enabled: true, Profile: "test"}

	testCases := []struct {
		name        string
		proxy       *ProxyConfig
		expectNil   bool
		expectValue *ProxyConfig
	}{
		{"with proxy", proxy, false, proxy},
		{"nil proxy", nil, true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{Proxy: tc.proxy}
			result := c.GetProxy()
			if tc.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectValue, result)
			}
		})
	}
}

func TestBaseScraperConfig_GetDownloadProxy(t *testing.T) {
	proxy := &ProxyConfig{Enabled: true, Profile: "dl-test"}

	testCases := []struct {
		name        string
		proxy       *ProxyConfig
		expectNil   bool
		expectValue *ProxyConfig
	}{
		{"with download proxy", proxy, false, proxy},
		{"nil download proxy", nil, true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := BaseScraperConfig{DownloadProxy: tc.proxy}
			result := c.GetDownloadProxy()
			if tc.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectValue, result)
			}
		})
	}
}

func TestBaseScraperConfig_EmbeddedSatisfiesInterface(t *testing.T) {
	type EmbeddedConfig struct {
		BaseScraperConfig `yaml:",inline"`
		Language          string `yaml:"language" json:"language"`
	}

	c := EmbeddedConfig{
		BaseScraperConfig: BaseScraperConfig{
			Enabled:      true,
			UserAgent:    "test-ua",
			RequestDelay: 100,
			MaxRetries:   3,
		},
		Language: "en",
	}

	var iface scraperutil.ScraperConfigInterface = &c

	assert.True(t, iface.IsEnabled())
	assert.Equal(t, "test-ua", iface.GetUserAgent())
	assert.Equal(t, 100, iface.GetRequestDelay())
	assert.Equal(t, 3, iface.GetMaxRetries())
}

func TestProxyAsConfig(t *testing.T) {
	proxy := &ProxyConfig{Enabled: true, Profile: "test"}

	testCases := []struct {
		name     string
		input    any
		expected *ProxyConfig
	}{
		{"nil returns nil", nil, nil},
		{"proxy config returns same pointer", proxy, proxy},
		{"wrong type returns nil", "not a proxy", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ProxyAsConfig(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
