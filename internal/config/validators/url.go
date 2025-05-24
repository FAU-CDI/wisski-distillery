//spellchecker:words validators
package validators

//spellchecker:words errors
import (
	"fmt"
	"net/url"

	"errors"
)

// URL represents a url.URL that is marshaled as a string representing the url.
type URL url.URL

func (u *URL) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

func (u *URL) String() string {
	if u == nil {
		return ""
	}
	return (*url.URL)(u).String()
}

func (u *URL) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	pu, err := url.Parse(string(text))
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	*u = URL(*pu)
	return nil
}

var errNotValidHTTPSURL = errors.New("not a valid https URL")

func ValidateHTTPSURL(url **URL, dflt string) error {
	if (*url).String() == "" {
		*url = new(URL)
		if err := (*url).UnmarshalText([]byte(dflt)); err != nil {
			return err
		}
	}
	if (*url).Scheme != "https" {
		return fmt.Errorf("%w: %q has scheme %q", errNotValidHTTPSURL, *url, (*url).Scheme)
	}
	return nil
}
