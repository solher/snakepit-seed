package validators

import (
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/solher/snakepit-seed/models"
)

type (
	UsersAdmin struct {
		users
	}
)

func NewUsersAdmin(l *logrus.Entry) *UsersAdmin {
	return &UsersAdmin{
		users: *newUsers(l),
	}
}

func (v *UsersAdmin) Create(users []models.User) ([]models.User, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.create(users)
}

func (v *UsersAdmin) Signin(cred *models.Credentials) (*models.Credentials, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.signin(cred)
}

func (v *UsersAdmin) Update(user *models.User) (*models.User, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.update(user)
}

func (v *UsersAdmin) UpdatePassword(pwd *models.Password) (*models.Password, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.updatePassword(pwd)
}

func (v *UsersAdmin) Output(users []models.User) []models.User {
	return v.output(users)
}
