package url

import (
	"errors"
	"net/url"
	"strings"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
)

func validateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return errors.Join(xerrors.ErrInvalidUrl, err)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return xerrors.ErrInvalidUrl
	}
	if parsed.Hostname() == "" {
		return xerrors.ErrInvalidUrl
	}

	return nil
}
