package interactors

import (
	"encoding/json"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"

	"git.wid.la/versatile/versatile-server/constants"
	"git.wid.la/versatile/versatile-server/models"
)

type (
	Sessions struct {
		snakepit.Interactor
		Repo HTTPSender
	}
)

func NewSessions(
	c *viper.Viper,
	l *logrus.Entry,
	r HTTPSender,
) *Sessions {
	return &Sessions{
		Interactor: *snakepit.NewInteractor(c, l),
		Repo:       r,
	}
}

func (i *Sessions) Create(session *models.Session) (*models.Session, error) {
	if err := i.Repo.Send(
		"",
		"POST",
		i.Constants.GetString(constants.AuthServerURL)+"/sessions",
		session,
		session,
	); err != nil {
		return nil, merry.Here(err)
	}

	return session, nil
}

func (i *Sessions) Delete(token string) (*models.Session, error) {
	session := &models.Session{}

	if err := i.Repo.Send(
		"",
		"DELETE",
		i.Constants.GetString(constants.AuthServerURL)+"/sessions/"+token,
		nil,
		session,
	); err != nil {
		return nil, merry.Here(err)
	}

	return session, nil
}

func (i *Sessions) DeleteCascade(wg *sync.WaitGroup, users []models.User) {
	ownerTokens := []string{}
	for _, user := range users {
		ownerTokens = append(ownerTokens, user.OwnerToken)
	}

	m, _ := json.Marshal(ownerTokens)

	if err := i.Repo.Send(
		"",
		"DELETE",
		i.Constants.GetString(constants.AuthServerURL)+"/sessions?ownerTokens="+string(m),
		nil,
		nil,
	); err != nil {
		panic(err)
	}

	wg.Done()
}
