package models

type AuthServerPayload struct {
	// The user associated with the session.
	User *User `json:"user,omitempty"`
	// The role name of the session.
	Role Role `json:"role,omitempty"`
}
