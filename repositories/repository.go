package repositories

import (
	"encoding/json"
	"time"

	"github.com/solher/snakepit-seed/errs"

	"gopkg.in/h2non/gentleman.v1"

	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

type (
	response struct {
		Description string `json:description,omitempty"`
	}

	DatabaseRunner interface {
		Run(q arangolite.Runnable) ([]byte, error)
	}

	Repository struct {
		snakepit.Repository
		DB     DatabaseRunner
		Client *gentleman.Client
	}
)

func NewRepository(
	c *viper.Viper,
	l *logrus.Entry,
	j *snakepit.JSON,
	db DatabaseRunner,
	cli *gentleman.Client,
) *Repository {
	return &Repository{
		Repository: *snakepit.NewRepository(c, l, j),
		DB:         db,
		Client:     cli,
	}
}

func (r *Repository) Run(q arangolite.Runnable, response interface{}) error {
	start := time.Now()
	raw, err := r.DB.Run(q)
	if err != nil {
		return merry.Here(err)
	}
	snakepit.LogTime(r.Logger, "Database requesting", start)

	if response == nil {
		return nil
	}

	if err := r.JSON.Unmarshal(r.Logger, "Database response", raw, response); err != nil {
		return merry.Here(err)
	}

	return nil
}

func (r *Repository) Send(authPayload, method, url string, body, response interface{}) error {
	req := r.Client.Request()
	req.AddHeader("Auth-Server-Payload", authPayload)
	req.Method(method)
	if method == "POST" {
		req.JSON(body)
	}
	req.URL(url)

	start := time.Now()
	res, err := req.Send()
	if err != nil {
		return merry.Here(err)
	}
	snakepit.LogTime(r.Logger, "HTTP requesting", start)

	if !res.Ok {
		errRes := &struct {
			Description string `json:description,omitempty"`
		}{}
		if err := json.Unmarshal(res.Bytes(), errRes); err != nil {
			return merry.Here(err)
		}

		if res.StatusCode == 403 || res.StatusCode == 404 {
			return merry.Here(errs.NotFound)
		}

		return merry.Errorf("From distant service: %s", errRes.Description)
	}

	if response == nil {
		return nil
	}

	if err := r.JSON.Unmarshal(r.Logger, "HTTP response", res.Bytes(), response); err != nil {
		return merry.Here(err)
	}

	return nil
}
