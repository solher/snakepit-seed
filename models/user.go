package models

type User struct {
	Document
	// The user first name.
	FirstName string `json:"firstName,omitempty"`
	// The user last name.
	LastName string `json:"lastName,omitempty"`
	// The user email.
	Email string `json:"email,omitempty"`
	// A unique identifier across the auth system.
	OwnerToken string `json:"ownerToken,omitempty"`
	// The user password.
	Password string `json:"password,omitempty"`
	// The role name of the user.
	Role Role `json:"role,omitempty"`
}

// swagger:response UsersResponse
type usersResponse struct {
	// in: body
	Body []User
}

// swagger:response UserResponse
type userResponse struct {
	// in: body
	Body User
}

// swagger:parameters UsersFindByKey UsersDeleteByKey UsersUpdateByKey UsersUpdatePassword
type usersKeyParam struct {
	// User key
	//
	// required: true
	// in: path
	Key string
}

// swagger:parameters UsersFind UsersFindByKey UsersDelete UsersUpdate
type usersFilterParam struct {
	// JSON filter defining offset, limit, sort, where and options
	//
	// in: query
	Filter string
}

// swagger:parameters UsersCreate UsersUpdate UsersUpdateByKey UsersSignup
type usersBodyParam struct {
	// required: true
	// in: body
	Body User
}
