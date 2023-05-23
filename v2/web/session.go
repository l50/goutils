package web

import "github.com/l50/goutils/v2/web/chrome"

// Credential contains the information that
// makes up a credential to authenticate
// to an application.
type Credential struct {
	User     string
	Password string
	TOTP     string
}

// Session contains parameters associated
// with maintaining a session.
type Session struct {
	Credential Credential
	Driver     *chrome.Driver
}
