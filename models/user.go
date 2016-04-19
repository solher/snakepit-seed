package models

type User struct {
	// The document handle. Format: ':collection/:key'
	ID string `json:"_id"`
	// The document's revision token. Changes at each update.
	Rev string `json:"_rev"`
	// The document's unique key.
	Key string `json:"_key"`
	// The user first name.
	FirstName string `json:"firstName"`
	// The user last name.
	LastName string `json:"lastName"`
	// The user email.
	Email string `json:"email"`
	// A unique identifier across the auth system.
	OwnerToken string `json:"ownerToken,omitempty"`
	// The user password.
	Password string `json:"password,omitempty"`
	// The role name of the user.
	Role string `json:"role"`
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

// swagger:parameters UsersFindByKey UsersDeleteByKey UsersUpdateByKey
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

// swagger:parameters UsersCreate UsersUpdate UsersUpdateByEmail
type usersBodyParam struct {
	// required: true
	// in: body
	Body User
}
