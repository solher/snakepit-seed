package models

type AuthServerPayload struct {
	// The user associated with the session.
	User *User `json:"user"`
	// The role name of the session.
	Role string `json:"role"`
}
