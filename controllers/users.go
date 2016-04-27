package controllers

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/solher/snakepit-seed/errs"
	"github.com/solher/snakepit-seed/models"

	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

type (
	UsersContext struct {
		AccessToken    string
		CurrentUser    *models.User
		CurrentSession *models.Session
		Key            string
		Filter         *filters.Filter
	}

	UsersInter interface {
		Create(userID string, users []models.User) ([]models.User, error)
		Find(userID string, f *filters.Filter) ([]models.User, error)
		Update(userID string, user *models.User, f *filters.Filter) ([]models.User, error)
		Delete(userID string, f *filters.Filter) ([]models.User, error)

		FindByKey(userID, id string, f *filters.Filter) (*models.User, error)
		UpdateByKey(userID, id string, user *models.User) (*models.User, error)
		DeleteByKey(userID, id string) (*models.User, error)

		Signin(cred *models.Credentials, agent string) (*models.Session, error)
		Signout(accessToken string) (*models.Session, error)
		UpdatePassword(userID, id, password string) (*models.User, error)
	}

	UsersValidator interface {
		Signin(cred *models.Credentials) (*models.Credentials, error)
		Create(users []models.User) ([]models.User, error)
		Update(user *models.User) (*models.User, error)
		UpdatePassword(pwd *models.Password) (*models.Password, error)
		Output(users []models.User) []models.User
	}

	SessionsOutputValidator interface {
		Output(users []models.Session) []models.Session
	}

	Users struct {
		snakepit.Controller
		Context           *UsersContext
		Inter             UsersInter
		Validator         UsersValidator
		SessionsValidator SessionsOutputValidator
	}
)

func NewUsers(
	c *viper.Viper,
	l *logrus.Entry,
	j *snakepit.JSON,
	ctx *UsersContext,
	i UsersInter,
	v UsersValidator,
	sv SessionsOutputValidator,
) *Users {
	return &Users{
		Controller:        *snakepit.NewController(c, l, j),
		Context:           ctx,
		Inter:             i,
		Validator:         v,
		SessionsValidator: sv,
	}
}

// Signin swagger:route POST /users/signin Users UsersSignin
//
// Sign in
//
// Signs in a user and returns a new session.
//
// Responses:
//  200: SessionResponse
func (c *Users) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, cred); !ok {
		return
	}

	cred, err := c.Validator.Signin(cred)
	if err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	session, err := c.Inter.Signin(cred, r.UserAgent())
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	session = &c.SessionsValidator.Output([]models.Session{*session})[0]

	c.JSON.Render(ctx, w, http.StatusCreated, session)
}

// CurrentSession swagger:route GET /users/me/session Users UsersCurrentSession
//
// Current session
//
// Returns the current session.
//
// Responses:
//  200: SessionResponse
func (c *Users) CurrentSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session := &c.SessionsValidator.Output([]models.Session{*c.Context.CurrentSession})[0]

	c.JSON.Render(ctx, w, http.StatusOK, session)
}

// Signout swagger:route POST /users/me/signout Users UsersSignout
//
// Sign out
//
// Signs out the current user.
//
// Responses:
//  200: SessionResponse
func (c *Users) Signout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session, err := c.Inter.Signout(c.Context.AccessToken)
	if err != nil {
		c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session = &c.SessionsValidator.Output([]models.Session{*session})[0]

	c.JSON.Render(ctx, w, http.StatusOK, session)
}

// Find swagger:route GET /users Users UsersFind
//
// Find
//
// Finds all the users matched by filter from the data source.
//
// Responses:
//  200: UsersResponse
func (c *Users) Find(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	users, err := c.Inter.Find(c.Context.CurrentUser.ID, c.Context.Filter)
	if err != nil {
		switch {
		case merry.Is(err, errs.InvalidFilter):
			c.JSON.RenderError(ctx, w, 422, errs.APIInvalidFilter, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	users = c.Validator.Output(users)

	c.JSON.Render(ctx, w, http.StatusOK, users)
}

// FindByKey swagger:route GET /users/{key} Users UsersFindByKey
//
// Find by key
//
// Finds a user by key from the data source.
//
// Responses:
//  200: UserResponse
func (c *Users) FindByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := c.Inter.FindByKey(c.Context.CurrentUser.ID, "users/"+c.Context.Key, c.Context.Filter)
	if err != nil {
		switch {
		case merry.Is(err, errs.InvalidFilter):
			c.JSON.RenderError(ctx, w, 422, errs.APIInvalidFilter, err)
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user = &c.Validator.Output([]models.User{*user})[0]

	c.JSON.Render(ctx, w, http.StatusOK, user)
}

// Create swagger:route POST /users Users UsersCreate
//
// Create
//
// Creates one or multiple users in the data source.
//
// Responses:
//  201: UserResponse
func (c *Users) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var users []models.User

	ok, bulk := c.JSON.UnmarshalBodyBulk(ctx, w, r.Body, &users)
	if !ok {
		return
	}

	users, err := c.Validator.Create(users)
	if err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	users, err = c.Inter.Create(c.Context.CurrentUser.ID, users)
	if err != nil {
		c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	users = c.Validator.Output(users)

	if bulk {
		c.JSON.Render(ctx, w, http.StatusCreated, users)
	} else {
		c.JSON.Render(ctx, w, http.StatusCreated, users[0])
	}
}

// Delete swagger:route DELETE /users Users UsersDelete
//
// Delete
//
// Deletes all the users matched by filter in the data source.
//
// Responses:
//  200: UsersResponse
func (c *Users) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	users, err := c.Inter.Delete(c.Context.CurrentUser.ID, c.Context.Filter)
	if err != nil {
		switch {
		case merry.Is(err, errs.InvalidFilter):
			c.JSON.RenderError(ctx, w, 422, errs.APIInvalidFilter, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	users = c.Validator.Output(users)

	c.JSON.Render(ctx, w, http.StatusOK, users)
}

// DeleteByKey swagger:route DELETE /users/{key} Users UsersDeleteByKey
//
// Delete by key
//
// Deletes a user by key in the data source.
//
// Responses:
//  200: UserResponse
func (c *Users) DeleteByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := c.Inter.DeleteByKey(c.Context.CurrentUser.ID, "users/"+c.Context.Key)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user = &c.Validator.Output([]models.User{*user})[0]

	c.JSON.Render(ctx, w, http.StatusOK, user)
}

// Update swagger:route PUT /users Users UsersUpdate
//
// Update
//
// Updates all the users matched by filter in the data source.
//
// Responses:
//  200: UserResponse
func (c *Users) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, user); !ok {
		return
	}

	user, err := c.Validator.Update(user)
	if err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	users, err := c.Inter.Update(c.Context.CurrentUser.ID, user, c.Context.Filter)
	if err != nil {
		switch {
		case merry.Is(err, errs.InvalidFilter):
			c.JSON.RenderError(ctx, w, 422, errs.APIInvalidFilter, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	users = c.Validator.Output(users)

	c.JSON.Render(ctx, w, http.StatusOK, users)
}

// UpdateByKey swagger:route PUT /users/{key} Users UsersUpdateByKey
//
// Update by key
//
// Updates a user by key in the data source.
//
// Responses:
//  200: UserResponse
func (c *Users) UpdateByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, user); !ok {
		return
	}

	user, err := c.Validator.Update(user)
	if err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	user, err = c.Inter.UpdateByKey(c.Context.CurrentUser.ID, "users/"+c.Context.Key, user)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user = &c.Validator.Output([]models.User{*user})[0]

	c.JSON.Render(ctx, w, http.StatusOK, user)
}

// UpdatePassword swagger:route POST /users/{key}/password Users UsersUpdatePassword
//
// Update password
//
// Updates the user password in the data source.
//
// Responses:
//  200: UserResponse
func (c *Users) UpdatePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pwd := &models.Password{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, pwd); !ok {
		return
	}

	pwd, err := c.Validator.UpdatePassword(pwd)
	if err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	user, err := c.Inter.UpdatePassword(c.Context.CurrentUser.ID, c.Context.Key, pwd.Password)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user = &c.Validator.Output([]models.User{*user})[0]

	c.JSON.Render(ctx, w, http.StatusOK, user)
}
