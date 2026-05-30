package httpclient

import (
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewScraperHTTPClient_BasicConfig(t *testing.T) {
	cfg := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	client, err := NewScraperHTTPClient(cfg, nil, config.FlareSolverrConfig{})
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewScraperHTTPClient_WithHeaders(t *testing.T) {
	cfg := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	client, err := NewScraperHTTPClient(cfg, nil, config.FlareSolverrConfig{},
		WithScraperHeaders(map[string]string{"X-Custom": "test"}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewScraperHTTPClient_WithCookies(t *testing.T) {
	cfg := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	client, err := NewScraperHTTPClient(cfg, nil, config.FlareSolverrConfig{},
		WithScraperCookies(map[string]string{"session": "abc123"}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewScraperHTTPClient_WithProxyProfile(t *testing.T) {
	cfg := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	client, err := NewScraperHTTPClient(cfg, nil, config.FlareSolverrConfig{},
		WithProxyProfile(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestInitScraperClient_BasicConfig(t *testing.T) {
	settings := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	result := InitScraperClient(settings, nil, config.FlareSolverrConfig{})
	assert.NotNil(t, result)
	assert.NotNil(t, result.Client)
	assert.False(t, result.ProxyEnabled)
}

func TestInitScraperClient_WithProxy(t *testing.T) {
	settings := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
		Proxy: &config.ProxyConfig{
			Enabled: true,
			Profile: "test",
			Profiles: map[string]config.ProxyProfile{
				"test": {URL: "http://proxy:8080"},
			},
		},
	}

	result := InitScraperClient(settings, nil, config.FlareSolverrConfig{})
	assert.NotNil(t, result)
	assert.NotNil(t, result.Client)
	assert.True(t, result.ProxyEnabled)
}

func TestInitScraperClient_WithGlobalProxy(t *testing.T) {
	settings := &config.ScraperSettings{
		Enabled:    true,
		Timeout:    30,
		RetryCount: 3,
	}

	globalProxy := &config.ProxyConfig{
		Enabled: true,
		Profile: "test",
		Profiles: map[string]config.ProxyProfile{
			"test": {URL: "http://global-proxy:8080"},
		},
	}

	result := InitScraperClient(settings, globalProxy, config.FlareSolverrConfig{})
	assert.NotNil(t, result)
	assert.NotNil(t, result.Client)
	assert.True(t, result.ProxyEnabled)
}
