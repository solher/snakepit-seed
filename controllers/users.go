package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"git.wid.la/versatile/versatile-server/errs"
	"git.wid.la/versatile/versatile-server/models"

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
		URLParams      map[string]string
		Filter         *filters.Filter
	}

	UsersInter interface {
		Find(userID string, f *filters.Filter) ([]models.User, error)
		FindByCred(cred *models.Credentials) (*models.User, error)
		FindByKey(userID, id string, f *filters.Filter) (*models.User, error)
		Create(userID string, users []models.User) ([]models.User, error)
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
		Signin(cred *models.Credentials) error
		Creation(users []models.User) error
		Update(user *models.User) error
	}

	Users struct {
		snakepit.Controller
		Context       UsersContext
		Inter         UsersInter
		SessionsInter SessionsReaderWriter
		Validator     UsersValidator
	}
)

func NewUsers(
	c *viper.Viper,
	l *logrus.Entry,
	j *snakepit.JSON,
	ctx UsersContext,
	i UsersInter,
	si SessionsReaderWriter,
	v UsersValidator,
) *Users {
	return &Users{
		Controller:    *snakepit.NewController(c, l, j),
		Context:       ctx,
		Inter:         i,
		SessionsInter: si,
		Validator:     v,
	}
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
func (c *Users) Signin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, cred); !ok {
		return
	}

	if err := c.Validator.Signin(cred); err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	user, err := c.Inter.FindByCred(cred)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	m, _ := json.Marshal(user)

	session := &models.Session{
		OwnerToken: user.OwnerToken,
		Agent:      r.UserAgent(),
		Policies:   []string{c.Constants.GetString(constants.PolicyName)},
		Payload:    string(m),
	}

	session, err = c.SessionsInter.Create(session)
	if err != nil {
		c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session.Policies = nil
	session.Payload = ""

	c.JSON.Render(ctx, w, http.StatusCreated, session)
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
func (c *Users) CurrentSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c.JSON.Render(ctx, w, http.StatusOK, c.Context.CurrentSession)
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
func (c *Users) Signout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session, err := c.SessionsInter.Delete(c.Context.AccessToken)
	if err != nil {
		c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	session.Policies = nil
	session.Payload = ""

	c.JSON.Render(ctx, w, http.StatusOK, session)
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

	for i := range users {
		users[i].Password = ""
	}

	c.JSON.Render(ctx, w, http.StatusOK, users)
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
func (c *Users) FindSelf(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := c.Inter.FindByKey(c.Context.CurrentUser.ID, c.Context.CurrentUser.ID, nil)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	c.JSON.Render(ctx, w, http.StatusOK, user)
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
func (c *Users) FindByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := c.Inter.FindByKey(c.Context.CurrentUser.ID, "users/"+c.Context.URLParams["key"], c.Context.Filter)
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

	user.Password = ""

	c.JSON.Render(ctx, w, http.StatusOK, user)
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
func (c *Users) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var users []models.User

	ok, bulk := c.JSON.UnmarshalBodyBulk(ctx, w, r.Body, &users)
	if !ok {
		return
	}

	if err := c.Validator.Creation(users); err != nil {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
		return
	}

	users, err := c.Inter.Create(c.Context.CurrentUser.ID, users)
	if err != nil {
		c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		return
	}

	if bulk {
		c.JSON.Render(ctx, w, http.StatusCreated, users)
	} else {
		c.JSON.Render(ctx, w, http.StatusCreated, users[0])
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

	c.JSON.Render(ctx, w, http.StatusOK, users)
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
func (c *Users) DeleteByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user, err := c.Inter.DeleteByKey(c.Context.CurrentUser.ID, "users/"+c.Context.URLParams["key"])
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.JSON.Render(ctx, w, http.StatusOK, user)
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
func (c *Users) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, user); !ok {
		return
	}
	// if err := c.Validator.UpdateOne(user); err != nil {
	// 	c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
	// 	return
	// }

	user.Key = ""
	user.Password = ""

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

	c.JSON.Render(ctx, w, http.StatusOK, users)
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
func (c *Users) UpdateByKey(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, user); !ok {
		return
	}

	// if err := c.Validator.UpdateOne(user); err != nil {
	// 	c.JSON.RenderError(ctx, w, 422, errs.APIValidation, err)
	// 	return
	// }

	user.Key = ""
	user.Password = ""

	user, err := c.Inter.UpdateByKey(c.Context.CurrentUser.ID, "users/"+c.Context.URLParams["key"], user)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	c.JSON.Render(ctx, w, http.StatusOK, user)
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
func (c *Users) UpdateSelfPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pwd := &models.Password{}

	if ok := c.JSON.UnmarshalBody(ctx, w, r.Body, pwd); !ok {
		return
	}

	if len(pwd.Password) == 0 {
		c.JSON.RenderError(ctx, w, 422, errs.APIValidation, merry.New("password cannot be blank"))
		return
	}

	user, err := c.Inter.UpdatePassword(c.Context.CurrentUser.ID, c.Context.CurrentUser.ID, pwd.Password)
	if err != nil {
		switch {
		case merry.Is(err, errs.NotFound):
			c.JSON.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
		default:
			c.JSON.RenderError(ctx, w, http.StatusInternalServerError, errs.APIInternal, err)
		}
		return
	}

	user.Password = ""

	c.JSON.Render(ctx, w, http.StatusOK, user)
}

// swagger:parameters Users UsersSignout UsersFindSelf
type tokenParam struct {
	// Access token (can also be set via the 'Authorization' header. Ex: 'Authorization: Bearer jhPd6Gf3jIP2h')
	//
	// in: query
	AccessToken string `json:"accessToken"`
}
