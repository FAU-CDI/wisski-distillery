//spellchecker:words validators
package validators

//spellchecker:words github errors
import (
	"net/url"

	"github.com/pkg/errors"
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
		return err
	}
	*u = URL(*pu)
	return nil
}

func ValidateHTTPSURL(url **URL, dflt string) error {
	if (*url).String() == "" {
		*url = new(URL)
		if err := (*url).UnmarshalText([]byte(dflt)); err != nil {
			return err
		}
	}
	if (*url).Scheme != "https" {
		return errors.Errorf("%q is not a valid https URL (%q)", *url, (*url).Scheme)
	}
	return nil
}
