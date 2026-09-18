package shortener

import (
	"errors"
)

func (s *Service) validateShortName(shortName string) error {
	if len(shortName) < s.opts.shortNameMinLength || len(shortName) > s.opts.shortNameMaxLength {
		return errors.New("short name must be between 3 and 32 characters long")
	}

	if !s.opts.shortNameRegexPattern.MatchString(shortName) {
		return errors.New("short name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}
