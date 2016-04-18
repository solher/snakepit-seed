package validators

import (
	"github.com/solher/snakepit"

	"git.wid.la/versatile/versatile-server/errs"
	"git.wid.la/versatile/versatile-server/models"
)

type (
	Users struct{}
)

func NewUsers() *Users {
	return &Users{}
}

func (v *Users) Creation(users []models.User) error {
	for _, user := range users {
		if len(user.Email) == 0 {
			return snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank)
		}

		if len(user.Password) == 0 {
			return snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank)
		}

		// if err := v.ValidateEmailUniqueness(user); err != nil {
		// 	return err
		// }
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
	if len(user.Email) == 0 {
		return snakepit.NewValidationError(errs.FieldEmail, errs.ValidBlank)
	}

	if len(user.Password) == 0 {
		return snakepit.NewValidationError(errs.FieldPassword, errs.ValidBlank)
	}

	return nil
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
