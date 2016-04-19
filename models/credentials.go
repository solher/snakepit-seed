package models

type Credentials struct {
	// The user email.
	Email string `json:"email"`
	// The user password.
	Password string `json:"password"`
}

// swagger:parameters UsersSignin
type credentialsBodyParam struct {
	// required: true
	// in: body
	Body Credentials
}

type Password struct {
	// The user password.
	Password string `json:"password"`
}

// swagger:parameters UsersUpdateSelfPassword
type passwordBodyParam struct {
	// required: true
	// in: body
	Body Password
}
