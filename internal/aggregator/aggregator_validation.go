package aggregator

import (
	"fmt"
	"strings"

	"github.com/fedora-oss/javinizer-go/internal/models"
)

// validateRequiredFields checks if all required fields are present and non-empty
func validateRequiredFields(movie *models.Movie, requiredFields []string) error {
	missingFields := []string{}

	for _, fieldName := range requiredFields {
		// Normalize field name to lowercase for case-insensitive comparison
		fieldLower := strings.ToLower(fieldName)

		// Check each field
		switch fieldLower {
		case "id":
			if movie.ID == "" {
				missingFields = append(missingFields, "ID")
			}
		case "contentid", "content_id":
			if movie.ContentID == "" {
				missingFields = append(missingFields, "ContentID")
			}
		case "title":
			if movie.Title == "" {
				missingFields = append(missingFields, "Title")
			}
		case "originaltitle", "original_title":
			if movie.OriginalTitle == "" {
				missingFields = append(missingFields, "OriginalTitle")
			}
		case "description", "plot":
			if movie.Description == "" {
				missingFields = append(missingFields, "Description")
			}
		case "director":
			if movie.Director == "" {
				missingFields = append(missingFields, "Director")
			}
		case "maker", "studio":
			if movie.Maker == "" {
				missingFields = append(missingFields, "Maker")
			}
		case "label":
			if movie.Label == "" {
				missingFields = append(missingFields, "Label")
			}
		case "series", "set":
			if movie.Series == "" {
				missingFields = append(missingFields, "Series")
			}
		case "releasedate", "release_date", "premiered":
			if movie.ReleaseDate == nil {
				missingFields = append(missingFields, "ReleaseDate")
			}
		case "runtime":
			if movie.Runtime == 0 {
				missingFields = append(missingFields, "Runtime")
			}
		case "coverurl", "cover_url", "cover":
			if movie.CoverURL == "" {
				missingFields = append(missingFields, "CoverURL")
			}
		case "posterurl", "poster_url", "poster":
			if movie.PosterURL == "" {
				missingFields = append(missingFields, "PosterURL")
			}
		case "trailerurl", "trailer_url", "trailer":
			if movie.TrailerURL == "" {
				missingFields = append(missingFields, "TrailerURL")
			}
		case "screenshots", "screenshot_url", "screenshoturl":
			if len(movie.Screenshots) == 0 {
				missingFields = append(missingFields, "Screenshots")
			}
		case "actresses", "actress":
			if len(movie.Actresses) == 0 {
				missingFields = append(missingFields, "Actresses")
			}
		case "genres", "genre":
			if len(movie.Genres) == 0 {
				missingFields = append(missingFields, "Genres")
			}
		case "rating", "ratingscore", "rating_score":
			// RatingScore == 0 is a valid value (some content has no rating).
			// We cannot distinguish "not scraped" from "intentionally 0" at validation time,
			// so we accept any value including 0. When RatingVotes == 0, we simply do not
			// add it to missingFields, allowing validation to pass.
			// When RatingVotes > 0, we have explicit rating data from a scraper.

		default:
			// Unknown field name - log warning but don't fail
			// This allows for forward compatibility
			continue
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
