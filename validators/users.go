package validators

import (
	"github.com/solher/snakepit"

	"github.com/solher/snakepit-seed/errs"
	"github.com/solher/snakepit-seed/middlewares"
	"github.com/solher/snakepit-seed/models"
)

type (
	Users struct{}
)

func NewUsers() *Users {
	return &Users{}
}

func (v *Users) Create(users []models.User) error {
	for _, user := range users {
		if len(user.Email) == 0 {
			return snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank)
		}

		if len(user.Password) == 0 {
			return snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank)
		}

		if len(user.Role) == 0 {
			return snakepit.NewValidationError(errs.FieldRole, errs.ValidBlank)
		}

		if err := v.roleExistence(middlewares.Role(user.Role)); err != nil {
			return err
		}
	}

	return nil
}

func (v *Users) Signin(cred *models.Credentials) error {
	if len(cred.Email) == 0 {
		return snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank)
	}

	if len(cred.Password) == 0 {
		return snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank)
	}

	return nil
}

func (v *Users) Update(user *models.User) error {
	if err := v.roleExistence(middlewares.Role(user.Role)); err != nil {
		return err
	}

	user.Key = ""
	user.Password = ""
	user.OwnerToken = ""

	return nil
}

func (v *Users) UpdatePassword(pwd *models.Password) error {
	if len(pwd.Password) == 0 {
		return snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank)
	}

	return nil
}

func (v *Users) Output(users []models.User) []models.User {
	for i := range users {
		users[i].Password = ""
	}

	return users
}

func (v *Users) roleExistence(role middlewares.Role) error {
	switch role {
	case middlewares.Admin,
		middlewares.Developer,
		middlewares.User:
		return nil
	}

	return snakepit.NewValidationError(errs.FieldRole, errs.ValidInvalid)
}

// func (v *Users) ValidateEmailUniqueness(user *models.User) error {
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
