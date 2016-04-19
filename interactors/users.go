package interactors

import (
	"sync"

	"github.com/solher/snakepit-seed/utils"

	"github.com/solher/snakepit-seed/errs"

	"golang.org/x/crypto/bcrypt"

	"github.com/solher/snakepit-seed/models"
	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/arangolite"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

type (
	SessionsCascadeDeleter interface {
		DeleteCascade(wg *sync.WaitGroup, users []models.User)
	}

	Users struct {
		snakepit.Interactor
		Repo          QueryRunner
		SessionsInter SessionsCascadeDeleter
	}
)

func NewUsers(
	c *viper.Viper,
	l *logrus.Entry,
	r QueryRunner,
	si SessionsCascadeDeleter,
) *Users {
	return &Users{
		Interactor:    *snakepit.NewInteractor(c, l),
		Repo:          r,
		SessionsInter: si,
	}
}

func (i *Users) Find(userID string, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, merry.Here(err)
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u._id == @userID
		%s
		RETURN u
	`, filter).Bind("userID", userID)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, merry.Here(err)
	}

	return users, nil
}

func (i *Users) FindByCred(cred *models.Credentials) (*models.User, error) {
	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.email == @email
		RETURN u
	`).Bind("email", cred.Email)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, merry.Here(err)
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

func (i *Users) FindByKey(userID, id string, f *filters.Filter) (*models.User, error) {
	if f == nil {
		f = &filters.Filter{}
	}

	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Find(userID, f)
	if err != nil {
		return nil, merry.Here(err)
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) Create(userID string, users []models.User) ([]models.User, error) {
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
		return nil, merry.Here(err)
	}

	return users, nil
}

func (i *Users) Delete(userID string, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, merry.Here(err)
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.createdBy == @userID
		%s
		REMOVE u IN users
		RETURN OLD
	`, filter).Bind("userID", userID)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, merry.Here(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go i.SessionsInter.DeleteCascade(wg, users)

	wg.Wait()

	return users, nil
}

func (i *Users) DeleteByKey(userID, id string) (*models.User, error) {
	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Delete(userID, f)
	if err != nil {
		return nil, merry.Here(err)
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) Update(userID string, user *models.User, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, merry.Here(err)
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u._id == @userID
		%s
		UPDATE u with @user IN users
		RETURN NEW
	`, filter).Bind("user", user).Bind("userID", userID)

	users := []models.User{}

	if err := i.Repo.Run(q, &users); err != nil {
		return nil, merry.Here(err)
	}

	return users, nil
}

func (i *Users) UpdateByKey(userID, id string, user *models.User) (*models.User, error) {
	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Update(userID, user, f)
	if err != nil {
		return nil, merry.Here(err)
	}

	if len(users) == 0 {
		return nil, merry.Here(errs.NotFound)
	}

	return &users[0], nil
}

func (i *Users) UpdatePassword(userID, id, password string) (*models.User, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return nil, merry.Here(err)
	}

	user := &models.User{
		Password: string(enc),
	}

	user, err = i.UpdateByKey(userID, id, user)
	if err != nil {
		return nil, merry.Here(err)
	}

	return user, nil
}
