package javstash

import (
	"fmt"
	"os"
	"strings"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/fedora-oss/javinizer-go/internal/configutil"
)

type JavstashConfig struct {
	config.BaseScraperConfig `yaml:",inline"`
	APIKey                   string `yaml:"api_key" json:"api_key"`
	BaseURL                  string `yaml:"base_url" json:"base_url"`
	Language                 string `yaml:"language" json:"language"`
}

func (c *JavstashConfig) ValidateConfig(sc *config.ScraperSettings) error {
	if err := config.ValidateCommonSettings("javstash", sc); err != nil {
		return err
	}
	if !sc.Enabled {
		return nil
	}
	apiKey := ""
	if v, ok := sc.Extra["api_key"].(string); ok {
		apiKey = strings.TrimSpace(v)
	}
	if apiKey == "" {
		apiKey = os.Getenv("JAVSTASH_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("javstash: api_key is required (set in config or JAVSTASH_API_KEY env var)")
	}
	if err := configutil.ValidateHTTPBaseURL("javstash.base_url", sc.BaseURL); err != nil {
		return err
	}
	return nil
}
