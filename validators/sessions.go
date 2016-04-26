package validators

import (
	"github.com/Sirupsen/logrus"
	"github.com/solher/snakepit"
	"github.com/solher/snakepit-seed/models"
)

type (
	Sessions struct {
		snakepit.Validator
	}
)

func NewSessions(l *logrus.Entry) *Sessions {
	return &Sessions{
		Validator: *snakepit.NewValidator(l),
	}
}

func (v *Sessions) Output(sessions []models.Session) []models.Session {
	for i := range sessions {
		sessions[i].Policies = nil
		sessions[i].Payload = ""
		sessions[i].OwnerToken = ""
	}

	return sessions
}
