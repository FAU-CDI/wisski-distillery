//spellchecker:words server
package server

// CSRFCookie, CSRFCookieField, SessionCookie and SessionUserKey
// hold the names of the cookies and fields used for specific cookies.
//
// These are intentionally kept short to conserve bandwidth.
const (
	CSRFCookie      = "F" // CSRF cookie sent on a lot of requests
	CSRFCookieField = "@" // form field name __should not be used by anything else__
	// to pay respect

	SessionCookie  = "x" // name of the cookie to use ; to doubt
	SessionUserKey = "@" // key within the session data to hold the username
)
