package repositories

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"git.wid.la/versatile/versatile-server/errs"

	"gopkg.in/h2non/gentleman-retry.v0"
	"gopkg.in/h2non/gentleman.v0"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/juju/ratelimit"
	"github.com/solher/arangolite"
)

type (
	response struct {
		Raw string `json:"raw,omitempty"`
	}

	DatabaseRunner interface {
		RunAsync(q arangolite.Runnable) (*arangolite.Result, error)
		Run(q arangolite.Runnable) ([]byte, error)
	}

	Repository struct {
		db DatabaseRunner
		c  *gentleman.Client
		b  *ratelimit.Bucket
	}
)

func New(db DatabaseRunner) *Repository {
	cli := gentleman.New()
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)))

	return &Repository{db: db, c: cli, b: ratelimit.NewBucket(10*time.Millisecond, 10)}
}

func (r *Repository) Run(q arangolite.Runnable) (*arangolite.Result, error) {
	r.b.Wait(1)

	result, err := r.db.RunAsync(q)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) RunSync(q arangolite.Runnable) ([]byte, error) {
	r.b.Wait(1)

	result, err := r.db.Run(q)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) Send(authPayload, method, url string, data interface{}) ([]byte, error) {
	req := r.c.Request()
	req.AddHeader("Auth-Server-Payload", authPayload)
	req.Method(method)
	if method == "POST" {
		req.JSON(data)
	}
	req.URL(url)

	r.b.Wait(1)

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	if !res.Ok {
		errRes := &response{}
		if err := json.Unmarshal(res.Bytes(), errRes); err != nil {
			return nil, err
		}

		if strings.Contains(errRes.Raw, "not found") {
			return nil, errs.NotFound
		}

		return nil, errors.New(errRes.Raw)
	}

	return res.Bytes(), nil
}
