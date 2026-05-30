package aventertainment

import (
	"fmt"
	"strings"

	"github.com/fedora-oss/javinizer-go/internal/config"
	"github.com/fedora-oss/javinizer-go/internal/configutil"
)

type AVEntertainmentConfig struct {
	config.BaseScraperConfig `yaml:",inline"`
	Language                 string `yaml:"language" json:"language"`
	BaseURL                  string `yaml:"base_url" json:"base_url"`
	ScrapeBonusScreens       bool   `yaml:"scrape_bonus_screens" json:"scrape_bonus_screens"`
}

func (c *AVEntertainmentConfig) ValidateConfig(sc *config.ScraperSettings) error {
	if err := config.ValidateCommonSettings("aventertainment", sc); err != nil {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(sc.Language)) {
	case "", "en", "ja":
	default:
		return fmt.Errorf("aventertainment: language must be 'en' or 'ja', got %q", sc.Language)
	}
	if err := configutil.ValidateHTTPBaseURL("aventertainment.base_url", sc.BaseURL); err != nil {
		return err
	}
	return nil
}
