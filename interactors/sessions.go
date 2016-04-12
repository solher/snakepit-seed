package interactors

import (
	"encoding/json"
	"sync"

	"git.wid.la/versatile/versatile-server/constants"
	"git.wid.la/versatile/versatile-server/models"
)

type (
	Sessions struct {
		g ConstantsGetter
		r HTTPSender
	}
)

func NewSessions(g ConstantsGetter, r HTTPSender) *Sessions {
	return &Sessions{g: g, r: r}
}

func (i *Sessions) Create(session *models.Session) (*models.Session, error) {
	res, err := i.r.Send("", "POST", i.g.GetString(constants.AuthServerURL)+"/sessions", session)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(res, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (i *Sessions) Delete(token string) (*models.Session, error) {
	res, err := i.r.Send("", "DELETE", i.g.GetString(constants.AuthServerURL)+"/sessions/"+token, nil)
	if err != nil {
		return nil, err
	}

	session := &models.Session{}

	if err := json.Unmarshal(res, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (i *Sessions) DeleteCascade(wg *sync.WaitGroup, users []models.User) {
	ownerTokens := []string{}
	for _, user := range users {
		ownerTokens = append(ownerTokens, user.OwnerToken)
	}

	m, _ := json.Marshal(ownerTokens)

	if _, err := i.r.Send("", "DELETE", i.g.GetString(constants.AuthServerURL)+"/sessions?ownerTokens="+string(m), nil); err != nil {
		panic(err)
	}

	wg.Done()
}
