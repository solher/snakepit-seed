package models

import "time"

type Session struct {
	// The creation timestamp.
	Created *time.Time `json:"created,omitempty"`
	// The validity time limit of the session.
	ValidTo *time.Time `json:"validTo,omitempty"`
	// The authentication token identifying the session.
	Token string `json:"token,omitempty"`
	// An optional token to find a user's sessions.
	OwnerToken string `json:"ownerToken,omitempty"`
	// The end user agent.
	// required: true
	Agent string `json:"agent,omitempty"`
	// The list of the policy names associated with the session.
	// required: true
	Policies []string `json:"policies,omitempty"`
	// A client non checked custom payload.
	Payload string `json:"payload,omitempty"`
	// The role name of the session.
	Role string `json:"role,omitempty"`
}

// swagger:response SessionResponse
type sessionResponse struct {
	// in: body
	Body Session
}
