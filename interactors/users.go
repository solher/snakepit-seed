package interactors

import (
	"encoding/json"
	"sync"

	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/middlewares"
	"github.com/solher/snakepit-seed/utils"

	"github.com/solher/snakepit-seed/errs"

	"golang.org/x/crypto/bcrypt"

	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/arangolite"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit"
	"github.com/solher/snakepit-seed/models"
	"github.com/spf13/viper"
)

type (
	SessionsReaderWriter interface {
		Create(session *models.Session) (*models.Session, error)
		Delete(token string) (*models.Session, error)
		DeleteCascade(wg *sync.WaitGroup, users []models.User)
	}

	Users struct {
		snakepit.Interactor
		Repo          QueryRunner
		SessionsInter SessionsReaderWriter
	}
)

func NewUsers(
	c *viper.Viper,
	l *logrus.Entry,
	r QueryRunner,
	si SessionsReaderWriter,
) *Users {
	return &Users{
		Interactor:    *snakepit.NewInteractor(c, l),
		Repo:          r,
		SessionsInter: si,
	}
}

func (i *Users) Signup(user *models.User) (*models.User, error) {
	user.Role = string(middlewares.User)

	users, err := i.Create([]models.User{*user})
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (i *Users) Find(f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, err
	}

	q := arangolite.NewQuery(`
		FOR u IN users
		%s
		RETURN u
	`, filter)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (i *Users) Signin(cred *models.Credentials, agent string) (*models.Session, error) {
	user, err := i.FindByCred(cred)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	payload := &models.AuthServerPayload{
		User: user,
		Role: user.Role,
	}

	m, _ := json.Marshal(payload)

	session := &models.Session{
		OwnerToken: user.OwnerToken,
		Agent:      agent,
		Policies:   []string{i.Constants.GetString(constants.PolicyName)},
		Payload:    string(m),
	}

	session, err = i.SessionsInter.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *Users) Signout(accessToken string) (*models.Session, error) {
	session, err := i.SessionsInter.Delete(accessToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *Users) FindByCred(cred *models.Credentials) (*models.User, error) {
	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.email == @email
		RETURN u
	`).Bind("email", cred.Email)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	user := &users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		return nil, merry.Here(errs.NotFound)
	}

	return user, nil
}

func (i *Users) FindByKey(key string, f *filters.Filter) (*models.User, error) {
	if f == nil {
		f = &filters.Filter{}
	}

	f.Where = append(f.Where, map[string]interface{}{"_key": key})

	users, err := i.Find(f)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) Create(users []models.User) ([]models.User, error) {
	for i := range users {
		enc, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), 11)
		if err != nil {
			return nil, merry.Here(err)
		}
		users[i].Password = string(enc)
		users[i].OwnerToken = utils.GenToken(32)
	}

	q := arangolite.NewQuery(`
		FOR u IN @users
		INSERT u IN users
		RETURN NEW
	`).Bind("users", users)

	users = []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (i *Users) Delete(f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, err
	}

	q := arangolite.NewQuery(`
		FOR u IN users
		%s
		REMOVE u IN users
		RETURN OLD
	`, filter)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go i.SessionsInter.DeleteCascade(wg, users)

	wg.Wait()

	return users, nil
}

func (i *Users) DeleteByKey(key string) (*models.User, error) {
	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_key": key})

	users, err := i.Delete(f)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) Update(user *models.User, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, err
	}

	q := arangolite.NewQuery(`
		FOR u IN users
		%s
		UPDATE u with @user IN users
		RETURN NEW
	`, filter).Bind("user", user)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (i *Users) UpdateByKey(key string, user *models.User) (*models.User, error) {
	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_key": key})

	users, err := i.Update(user, f)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) UpdatePassword(key, password string) (*models.User, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return nil, merry.Here(err)
	}

	user := &models.User{
		Password: string(enc),
	}

	user, err = i.UpdateByKey(key, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
