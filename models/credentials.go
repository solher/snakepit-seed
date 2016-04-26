package models

type Credentials struct {
	// The user email.
	Email string `json:"email,omitempty"`
	// The user password.
	Password string `json:"password,omitempty"`
}

// swagger:parameters UsersSignin
type credentialsBodyParam struct {
	// required: true
	// in: body
	Body Credentials
}

type Password struct {
	// The user password.
	Password string `json:"password,omitempty"`
}

// swagger:parameters UsersUpdatePassword
type passwordBodyParam struct {
	// required: true
	// in: body
	Body Password
}
