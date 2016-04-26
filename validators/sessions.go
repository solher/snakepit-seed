package validators

import "github.com/solher/snakepit-seed/models"

type (
	Sessions struct{}
)

func NewSessions() *Sessions {
	return &Sessions{}
}

func (v *Sessions) Output(sessions []models.Session) []models.Session {
	for i := range sessions {
		sessions[i].Policies = nil
		sessions[i].Payload = ""
		sessions[i].OwnerToken = ""
	}

	return sessions
}
