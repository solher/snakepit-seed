package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/pressly/chi"
	"git.wid.la/versatile/versatile-server/errs"
	"git.wid.la/versatile/versatile-server/middlewares"
	"git.wid.la/versatile/versatile-server/models"

	"github.com/palantir/stacktrace"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit"
)

type (
	UsersInter interface {
		Find(userID string, f *filters.Filter) ([]models.User, error)
		FindByCred(cred *models.Credentials) (*models.User, error)
		FindByKey(userID, id string, f *filters.Filter) (*models.User, error)
		Create(userID string, users []models.User) ([]models.User, error)
		CreateOne(userID string, user *models.User) (*models.User, error)
		Delete(userID string, f *filters.Filter) ([]models.User, error)
		DeleteByKey(userID, id string) (*models.User, error)
		Update(userID string, user *models.User, f *filters.Filter) ([]models.User, error)
		UpdateByKey(userID, id string, user *models.User) (*models.User, error)
		UpdatePassword(userID, id, password string) (*models.User, error)
	}

	SessionsReaderWriter interface {
		Create(session *models.Session) (*models.Session, error)
		Delete(token string) (*models.Session, error)
	}

	UsersValidator interface {
		ValidateSignin(cred *models.Credentials) error
		ValidateCreation(users []models.User) error
		ValidateUpdate(user *models.User) error
	}

	UsersCtrl struct {
		i   UsersInter
		srw SessionsReaderWriter
		v   UsersValidator
		r   *snakepit.Render
	}
)

func NewUsersCtrl(
	i UsersInter,
	srw SessionsReaderWriter,
	v UsersValidator,
	r *snakepit.Render,
) *UsersCtrl {
	return &UsersCtrl{i: i, srw: srw, v: v, r: r}
}

// Signin swagger:route POST /signin Users UsersSignin
//
// Sign in
//
// Signs in a user and returns a new session.
//
// Responses:
//  200: SessionResponse
//  400: BodyDecodingResponse
//  401: UnauthorizedResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *UsersCtrl) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}

	if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIBodyDecoding, err)
		return
	}

	if err := c.v.ValidateSignin(cred); err != nil {
		c.r.JSONError(w, 422, errs.APIValidation, err)
		return
	}

	user, err := c.i.FindByCred(cred)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	m, _ := json.Marshal(user)

	session := &models.Session{
		OwnerToken: user.OwnerToken,
		Agent:      r.UserAgent(),
		Policies:   []string{"co-net"},
		Payload:    string(m),
	}

	session, err = c.srw.Create(session)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session.Policies = nil
	session.Payload = ""

	c.r.JSON(w, http.StatusCreated, session)
}

// CurrentSession swagger:route GET /me/session Users UsersCurrentSession
//
// Current session
//
// Returns the current session.
//
// Responses:
//  200: SessionResponse
//  401: UnauthorizedResponse
//  500: InternalResponse
func (c *UsersCtrl) CurrentSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session, err := middlewares.SessionFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	c.r.JSON(w, http.StatusOK, session)
}

// Signout swagger:route POST /me/signout Users UsersSignout
//
// Sign out
//
// Signs out the current user.
//
// Responses:
//  200: SessionResponse
//  401: UnauthorizedResponse
//  500: InternalResponse
func (c *UsersCtrl) Signout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	token, err := middlewares.AccessTokenFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session, err := c.srw.Delete(token)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session.Policies = nil
	session.Payload = ""

	c.r.JSON(w, http.StatusOK, session)
}

// Find swagger:route GET / Users UsersFind
//
// Find
//
// Finds all the users matched by filter from the data source.
//
// Responses:
//  200: UsersResponse
//  400: FilterDecodingResponse
//  401: UnauthorizedResponse
//  422: InvalidFilterResponse
//  500: InternalResponse
func (c *UsersCtrl) Find(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	filter, err := filters.FromRequest(r)
	if err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIFilterDecoding, err)
		return
	}

	users, err := c.i.Find(currentUser.ID, filter)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeInvalidFilter:
			c.r.JSONError(w, 422, errs.APIInvalidFilter, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	c.r.JSON(w, http.StatusOK, users)
}

// FindSelf swagger:route GET /me Users UsersFindSelf
//
// Find self
//
// Finds the user currently signed in.
//
// Responses:
//  200: UserResponse
//  401: UnauthorizedResponse
func (c *UsersCtrl) FindSelf(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	user, err = c.i.FindByKey(user.ID, user.ID, nil)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	c.r.JSON(w, http.StatusOK, user)
}

// FindByKey swagger:route GET /{key} Users UsersFindByKey
//
// Find by key
//
// Finds a user by key from the data source.
//
// Responses:
//  200: UserResponse
//  400: FilterDecodingResponse
//  401: UnauthorizedResponse
//  422: InvalidFilterResponse
//  500: InternalResponse
func (c *UsersCtrl) FindByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	filter, err := filters.FromRequest(r)
	if err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIFilterDecoding, err)
		return
	}

	user, err := c.i.FindByKey(currentUser.ID, "users/"+chi.URLParams(ctx)["key"], filter)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeInvalidFilter:
			c.r.JSONError(w, 422, errs.APIInvalidFilter, err)
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	c.r.JSON(w, http.StatusOK, user)
}

// Create swagger:route POST / Users UsersCreate
//
// Create
//
// Creates one or multiple users in the data source.
//
// Responses:
//  201: UserResponse
//  400: BodyDecodingResponse
//  401: UnauthorizedResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *UsersCtrl) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	user := &models.User{}
	var users []models.User

	buffer, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(buffer, user); err != nil {
		if err := json.Unmarshal(buffer, &users); err != nil {
			c.r.JSONError(w, http.StatusBadRequest, errs.APIBodyDecoding, err)
			return
		}
	}

	if users == nil {
		err = c.v.ValidateCreation([]models.User{*user})
	} else {
		err = c.v.ValidateCreation(users)
	}

	if err != nil {
		c.r.JSONError(w, 422, errs.APIValidation, err)
		return
	}

	if users == nil {
		user, err = c.i.CreateOne(currentUser.ID, user)
	} else {
		users, err = c.i.Create(currentUser.ID, users)
	}

	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	if users == nil {
		c.r.JSON(w, http.StatusCreated, user)
	} else {
		c.r.JSON(w, http.StatusCreated, users)
	}
}

// Delete swagger:route DELETE / Users UsersDelete
//
// Delete
//
// Deletes all the users matched by filter in the data source.
//
// Responses:
//  200: UsersResponse
//  400: FilterDecodingResponse
//  401: UnauthorizedResponse
//  422: InvalidFilterResponse
//  500: InternalResponse
func (c *UsersCtrl) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	filter, err := filters.FromRequest(r)
	if err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIFilterDecoding, err)
		return
	}

	users, err := c.i.Delete(currentUser.ID, filter)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeInvalidFilter:
			c.r.JSONError(w, 422, errs.APIInvalidFilter, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, users)
}

// DeleteByKey swagger:route DELETE /{key} Users UsersDeleteByKey
//
// Delete by key
//
// Deletes a user by key in the data source.
//
// Responses:
//  200: UserResponse
//  401: UnauthorizedResponse
//  500: InternalResponse
func (c *UsersCtrl) DeleteByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	user, err := c.i.DeleteByKey(currentUser.ID, "users/"+chi.URLParams(ctx)["key"])
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, user)
}

// Update swagger:route PUT / Users UsersUpdate
//
// Update
//
// Updates all the users matched by filter in the data source.
//
// Responses:
//  200: UserResponse
//  401: UnauthorizedResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  400: FilterDecodingResponse
//  422: InvalidFilterResponse
//  500: InternalResponse
func (c *UsersCtrl) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	filter, err := filters.FromRequest(r)
	if err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIFilterDecoding, err)
		return
	}

	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIBodyDecoding, err)
		return
	}

	// if err := c.v.ValidateUpdateOne(user); err != nil {
	// 	c.r.JSONError(w, 422, errs.APIValidation, err)
	// 	return
	// }

	user.Key = ""
	user.Password = ""

	users, err := c.i.Update(currentUser.ID, user, filter)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeInvalidFilter:
			c.r.JSONError(w, 422, errs.APIInvalidFilter, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, users)
}

// UpdateByKey swagger:route PUT /{key} Users UsersUpdateByKey
//
// Update by key
//
// Updates a user by key in the data source.
//
// Responses:
//  200: UserResponse
//  401: UnauthorizedResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *UsersCtrl) UpdateByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	user := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIBodyDecoding, err)
		return
	}

	// if err := c.v.ValidateUpdateOne(user); err != nil {
	// 	c.r.JSONError(w, 422, errs.APIValidation, err)
	// 	return
	// }

	user.Key = ""
	user.Password = ""

	user, err = c.i.UpdateByKey(currentUser.ID, "users/"+chi.URLParams(ctx)["key"], user)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.r.JSON(w, http.StatusOK, user)
}

// UpdateSelfPassword swagger:route POST /me/password Users UsersUpdateSelfPassword
//
// Update self password
//
// Updates the current user password.
//
// Responses:
//  200: UserResponse
//  401: UnauthorizedResponse
//  400: BodyDecodingResponse
//  422: ValidationResponse
//  500: InternalResponse
func (c *UsersCtrl) UpdateSelfPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	currentUser, err := middlewares.UserFromCtx(ctx)
	if err != nil {
		c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	pwd := &models.Password{}

	if err := json.NewDecoder(r.Body).Decode(pwd); err != nil {
		c.r.JSONError(w, http.StatusBadRequest, errs.APIBodyDecoding, err)
		return
	}

	if len(pwd.Password) == 0 {
		c.r.JSONError(w, 422, errs.APIValidation, errors.New("password cannot be blank"))
		return
	}

	user, err := c.i.UpdatePassword(currentUser.ID, currentUser.ID, pwd.Password)
	if err != nil {
		switch stacktrace.GetCode(err) {
		case errs.EcodeNotFound:
			c.r.JSONError(w, http.StatusUnauthorized, errs.APIUnauthorized, err)
		default:
			c.r.JSONError(w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	c.r.JSON(w, http.StatusOK, user)
}

// swagger:parameters Users UsersSignout UsersFindSelf
type tokenParam struct {
	// Access token (can also be set via the 'Authorization' header. Ex: 'Authorization: Bearer jhPd6Gf3jIP2h')
	//
	// in: query
	AccessToken string `json:"accessToken"`
}
