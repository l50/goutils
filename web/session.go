package web

// Credential contains the information that
// makes up a credential to authenticate
// to an application.
type Credential struct {
	User       string
	Password   string
	TwoFacCode string
}

// FormField contains a form field name and its associated selector.
type FormField struct {
	Name     string `json:"-"`
	Selector string `json:"-"`
}

// Session contains parameters associated
// with maintaining a session.
type Session struct {
	Credential Credential
	Driver     interface{}
}
