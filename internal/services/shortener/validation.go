package shortener

import (
	"errors"
	"fmt"
	"net/url"
)

func (s *Service) validateOriginalURL(originalURL string) error {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return errors.New("invalid original URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("original URL must be an HTTP or HTTPS URL")
	}

	if parsedURL.Hostname() == "" {
		return errors.New("original URL must have a valid host")
	}

	return nil
}

func (s *Service) validateShortName(shortName string) error {
	if len(shortName) < s.opts.shortNameMinLength || len(shortName) > s.opts.shortNameMaxLength {
		return fmt.Errorf("short name must be between %d and %d characters long", s.opts.shortNameMinLength, s.opts.shortNameMaxLength)
	}

	if !s.opts.shortNameRegexPattern.MatchString(shortName) {
		return fmt.Errorf("short name must match the pattern: %s", s.opts.shortNameRegexPattern.String())
	}

	return nil
}
