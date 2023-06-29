package pwmgr

// Record represents a record in a password manager.
//
// **Attributes:**
//
// UID: A unique identifier.
// Title: Title of the record.
// URL: The associated URL of the record.
// Username: The username associated with the record.
// Password: The password associated with the record.
// TOTP: Time-based One-Time Password.
// Note: Additional note associated with the record.
type Record struct {
	UID      string
	Title    string
	URL      string
	Username string
	Password string
	TOTP     string
	Note     string
}

// PasswordManager represents a password manager.
type PasswordManager interface {
	IsInstalled() bool
	IsLoggedIn() bool
	RetrieveRecord(uid string) (Record, error)
	SearchRecords(searchTerm string) (string, error)
	AddRecord(fields map[string]string) error
}
