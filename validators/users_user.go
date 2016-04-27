package validators

import (
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/solher/snakepit-seed/models"
)

type (
	UsersUser struct {
		users
	}
)

func NewUsersUser(l *logrus.Entry) *UsersUser {
	return &UsersUser{
		users: *newUsers(l),
	}
}

func (v *UsersUser) Signup(user *models.User) (*models.User, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.signup(user)
}

// Create is never used by a user
func (v *UsersUser) Create(users []models.User) ([]models.User, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.create(users)
}

func (v *UsersUser) Signin(cred *models.Credentials) (*models.Credentials, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.signin(cred)
}

func (v *UsersUser) Update(user *models.User) (*models.User, error) {
	start := time.Now()
	defer v.LogTime(start)

	user.Role = ""

	return v.update(user)
}

func (v *UsersUser) UpdatePassword(pwd *models.Password) (*models.Password, error) {
	start := time.Now()
	defer v.LogTime(start)

	return v.updatePassword(pwd)
}

func (v *UsersUser) Output(users []models.User) []models.User {
	for i := range users {
		users[i].OwnerToken = ""
	}

	return v.output(users)
}
