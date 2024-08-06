package tgbot

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseURL(input string) (string, string, error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return "", "", fmt.Errorf("invalid protocol: %s", parsedURL.Scheme)
	}

	if parsedURL.Host != "collectible.money" {
		return "", "", fmt.Errorf("invalid hostname: %s", parsedURL.Host)
	}

	// Strip the query parameters
	parsedURL.RawQuery = ""

	// Split the path into segments
	segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	// Check if the URL has enough segments
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid URL data")
	}

	urlType := segments[0]
	slug := segments[1]

	validTypes := map[string]bool{"item": true, "series": true, "creator": true}
	if !validTypes[urlType] {
		return "", "", fmt.Errorf("invalid URL type: %s", urlType)
	}

	return urlType, slug, nil
}
