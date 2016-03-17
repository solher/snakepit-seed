package interactors

import (
	"encoding/json"
	"sync"

	"github.com/palantir/stacktrace"
	"git.wid.la/versatile/versatile-server/constants"
	"git.wid.la/versatile/versatile-server/models"
)

type (
	SessionsInter struct {
		g ConstantsGetter
		r HTTPSender
	}
)

func NewSessionsInter(g ConstantsGetter, r HTTPSender) *SessionsInter {
	return &SessionsInter{g: g, r: r}
}

func (i *SessionsInter) Create(session *models.Session) (*models.Session, error) {
	res, err := i.r.Send("", "POST", i.g.GetString(constants.AuthServerURL)+"/sessions", session)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	if err := json.Unmarshal(res, session); err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return session, nil
}

func (i *SessionsInter) Delete(token string) (*models.Session, error) {
	res, err := i.r.Send("", "DELETE", i.g.GetString(constants.AuthServerURL)+"/sessions/"+token, nil)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	session := &models.Session{}

	if err := json.Unmarshal(res, session); err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return session, nil
}

func (i *SessionsInter) DeleteCascade(wg *sync.WaitGroup, users []models.User) {
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
