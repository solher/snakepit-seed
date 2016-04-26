package validators

import (
	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/snakepit"

	"github.com/solher/snakepit-seed/errs"
	"github.com/solher/snakepit-seed/middlewares"
	"github.com/solher/snakepit-seed/models"
)

type (
	users struct {
		snakepit.Validator
	}
)

func newUsers(l *logrus.Entry) *users {
	return &users{
		Validator: *snakepit.NewValidator(l),
	}
}

func (v *users) create(users []models.User) ([]models.User, error) {
	for _, user := range users {
		if len(user.Email) == 0 {
			return nil, merry.Here(snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank))
		}

		if len(user.Password) == 0 {
			return nil, merry.Here(snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank))
		}

		if len(user.Role) == 0 {
			return nil, merry.Here(snakepit.NewValidationError(errs.FieldRole, errs.ValidBlank))
		}

		if err := v.roleExistence(middlewares.Role(user.Role)); err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (v *users) signin(cred *models.Credentials) (*models.Credentials, error) {
	if len(cred.Email) == 0 {
		return nil, merry.Here(snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank))
	}

	if len(cred.Password) == 0 {
		return nil, merry.Here(snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank))
	}

	return cred, nil
}

func (v *users) update(user *models.User) (*models.User, error) {
	if err := v.roleExistence(middlewares.Role(user.Role)); err != nil {
		return nil, err
	}

	user.Key = ""
	user.Password = ""
	user.OwnerToken = ""

	return user, nil
}

func (v *users) updatePassword(pwd *models.Password) (*models.Password, error) {
	if len(pwd.Password) == 0 {
		return nil, merry.Here(snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank))
	}

	return pwd, nil
}

func (v *users) output(users []models.User) []models.User {
	for i := range users {
		users[i].Password = ""
	}

	return users
}

func (v *users) roleExistence(role middlewares.Role) error {
	switch role {
	case middlewares.Admin,
		middlewares.Developer,
		middlewares.User:
		return nil
	case "":
		return nil
	}

	return merry.Here(snakepit.NewValidationError(errs.FieldRole, errs.ValidInvalid))
}

// func (v *users) ValidateEmailUniqueness(user *models.User) error {
// 	if user.Email == nil {
// 		return nil
// 	}

// 	err := v.r.View(func(tx *bolt.Tx) error {
// 		raw := tx.Bucket([]byte("users")).Get([]byte(*user.Email))

// 		if len(raw) != 0 {
// 			return errs.NewErrValidation("email must be unique")
// 		}

// 		return nil
// 	})

// 	return err
// }
