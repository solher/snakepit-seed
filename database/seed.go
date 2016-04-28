package database

import (
	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/utils"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/solher/snakepit-seed/models"
)

type ProdSeed struct {
	Users []models.User `check:"keyOnly"`
}

func NewEmptyProdSeed() *ProdSeed {
	return &ProdSeed{}
}

func NewProdSeed() *ProdSeed {
	s := NewEmptyProdSeed()

	enc, _ := bcrypt.GenerateFromPassword([]byte("admin"), 11)

	s.Users = append(s.Users, []models.User{
		{
			Document:   models.NewDocument("", "", "admin"),
			FirstName:  "admin",
			LastName:   "admin",
			Email:      "admin",
			OwnerToken: utils.GenToken(32),
			Password:   string(enc),
			Role:       constants.RoleAdmin,
		},
	}...)

	return s
}

func (s *ProdSeed) PopulateConstants(v *viper.Viper) {
}
