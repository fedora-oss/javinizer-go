package system

import (
	"github.com/fedora-oss/javinizer-go/internal/scraperutil"
)

func scraperDisplayTitleAndOptions(name string, profileChoices []ScraperChoice) (string, []ScraperOption) {
	// Try to get self-registered options from scraper
	if provider, exists := scraperutil.GetScraperOptions(name); exists {
		// Convert []any to []ScraperOption
		options := make([]ScraperOption, 0, len(provider.Options)+10)
		for _, opt := range provider.Options {
			if so, ok := opt.(ScraperOption); ok {
				options = append(options, so)
			}
		}

		// Append common options (user agent, proxy, download proxy)
		options = append(options, scraperUserAgentOptions()...)
		options = append(options, scraperProxyOptions(profileChoices)...)
		options = append(options, scraperDownloadProxyOptions(profileChoices)...)

		return provider.DisplayTitle, options
	}

	// Fallback for scrapers without registered options:
	// Return defaults with just common options
	options := scraperUserAgentOptions()
	options = append(options, scraperProxyOptions(profileChoices)...)
	options = append(options, scraperDownloadProxyOptions(profileChoices)...)

	return name, options
}
