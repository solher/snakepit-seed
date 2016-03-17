package interactors

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/palantir/stacktrace"
	"git.wid.la/versatile/versatile-server/utils"

	"git.wid.la/versatile/versatile-server/errs"

	"gopkg.in/tomb.v2"

	"golang.org/x/crypto/bcrypt"

	"github.com/solher/arangolite"
	"github.com/solher/arangolite/filters"
	"git.wid.la/versatile/versatile-server/models"
)

type (
	GraphCascadeDeleter interface {
		DeleteCascade(wg *sync.WaitGroup, users []models.User)
	}

	SessionsCascadeDeleter interface {
		DeleteCascade(wg *sync.WaitGroup, users []models.User)
	}

	UsersInter struct {
		r    QueryRunner
		gcd  GraphCascadeDeleter
		scd  SessionsCascadeDeleter
		tPwd tomb.Tomb
	}
)

func NewUsersInter(r QueryRunner, gcd GraphCascadeDeleter, scd SessionsCascadeDeleter) *UsersInter {
	inter := &UsersInter{r: r, gcd: gcd, scd: scd}
	inter.tPwd.Kill(nil)
	return inter
}

func (i *UsersInter) Find(userID string, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.createdBy == @userID || u._id == @userID
		%s
		RETURN u
	`, filter).Bind("userID", userID)

	r, err := i.r.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	users := []models.User{}

	decoder := json.NewDecoder(r.Buffer())
	for r.HasMore() {
		batch := []models.User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	return users, nil
}

func (i *UsersInter) FindByCred(cred *models.Credentials) (*models.User, error) {
	if cred == nil {
		return nil, stacktrace.NewError("nil cred")
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.email == @email
		RETURN u
	`).Bind("email", cred.Email)

	r, err := i.r.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	users := []models.User{}

	decoder := json.NewDecoder(r.Buffer())
	for r.HasMore() {
		batch := []models.User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	if len(users) == 0 {
		return nil, stacktrace.Propagate(errs.ErrNotFound, "")
	}

	user := &users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		return nil, stacktrace.PropagateWithCode(err, errs.EcodeNotFound, "Wrong password")
	}

	return user, nil
}

func (i *UsersInter) FindByKey(userID, id string, f *filters.Filter) (*models.User, error) {
	if f == nil {
		f = &filters.Filter{}
	}

	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Find(userID, f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if len(users) == 0 {
		return nil, stacktrace.Propagate(errs.ErrNotFound, "")
	}

	return &users[0], nil
}

func (i *UsersInter) Create(userID string, users []models.User) ([]models.User, error) {
	if users == nil {
		return nil, stacktrace.NewError("nil users")
	}

	toGenerate := []models.User{}

	for i := range users {
		if users[i].Key != "" {
			user := models.User{}
			user.Key = users[i].Key
			user.Password = users[i].Password
			users[i].Password = ""
			toGenerate = append(toGenerate, user)
		} else {
			enc, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), 9)
			if err != nil {
				return nil, stacktrace.Propagate(err, "")
			}

			users[i].Password = string(enc)
		}

		users[i].OwnerToken = utils.GenToken(32)
	}

	q := arangolite.NewQuery(`
		FOR u IN @users
		INSERT u IN users
		RETURN NEW
	`).Bind("users", users)

	r, err := i.r.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	users = []models.User{}

	decoder := json.NewDecoder(r.Buffer())
	for r.HasMore() {
		batch := []models.User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	if i.tPwd.Alive() {
		i.tPwd.Kill(nil)
		_ = i.tPwd.Wait()
	}

	i.tPwd = tomb.Tomb{}

	i.tPwd.Go(func() error {
		bulk := []models.User{}

		for _, user := range toGenerate {
			select {
			case <-i.tPwd.Dying():
				return nil
			default:
			}

			time.Sleep(100 * time.Millisecond)

			enc, err := bcrypt.GenerateFromPassword([]byte(user.Password), 9)
			if err != nil {
				return err
			}

			user.Password = string(enc)
			bulk = append(bulk, user)

			if len(bulk) == 10 {
				_, _ = i.r.Run(arangolite.NewQuery(`
                    FOR u IN @users
                    UPDATE u IN users
                `).Bind("users", bulk))

				bulk = []models.User{}
			}
		}

		_, _ = i.r.Run(arangolite.NewQuery(`
            FOR u IN @users
            UPDATE u IN users
        `).Bind("users", bulk))

		return nil
	})

	return users, nil
}

func (i *UsersInter) CreateOne(userID string, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, stacktrace.NewError("nil user")
	}

	users, err := i.Create(userID, []models.User{*user})
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return &users[0], nil
}

func (i *UsersInter) Delete(userID string, f *filters.Filter) ([]models.User, error) {
	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.createdBy == @userID
		%s
		REMOVE u IN users
		RETURN OLD
	`, filter).Bind("userID", userID)

	r, err := i.r.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	users := []models.User{}

	decoder := json.NewDecoder(r.Buffer())
	for r.HasMore() {
		batch := []models.User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go i.gcd.DeleteCascade(wg, users)
	go i.scd.DeleteCascade(wg, users)

	wg.Wait()

	return users, nil
}

func (i *UsersInter) DeleteByKey(userID, id string) (*models.User, error) {
	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Delete(userID, f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if len(users) == 0 {
		return nil, stacktrace.Propagate(errs.ErrNotFound, "")
	}

	return &users[0], nil
}

func (i *UsersInter) Update(userID string, user *models.User, f *filters.Filter) ([]models.User, error) {
	if user == nil {
		return nil, stacktrace.NewError("nil user")
	}

	filter, err := utils.FilterToAQL("u", f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	q := arangolite.NewQuery(`
		FOR u IN users
        FILTER u.createdBy == @userID || u._id == @userID
		%s
		UPDATE u with @user IN users
		RETURN NEW
	`, filter).Bind("user", user).Bind("userID", userID)

	r, err := i.r.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	users := []models.User{}

	decoder := json.NewDecoder(r.Buffer())
	for r.HasMore() {
		batch := []models.User{}
		decoder.Decode(&batch)
		users = append(users, batch...)
	}

	return users, nil
}

func (i *UsersInter) UpdateByKey(userID, id string, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, stacktrace.NewError("nil user")
	}

	f := &filters.Filter{}
	f.Where = append(f.Where, map[string]interface{}{"_id": id})

	users, err := i.Update(userID, user, f)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if len(users) == 0 {
		return nil, stacktrace.Propagate(errs.ErrNotFound, "")
	}

	return &users[0], nil
}

func (i *UsersInter) UpdatePassword(userID, id, password string) (*models.User, error) {
	user := &models.User{}

	enc, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	user.Password = string(enc)

	user, err = i.UpdateByKey(userID, id, user)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return user, nil
}
