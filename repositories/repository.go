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
	"github.com/palantir/stacktrace"
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

func NewRepository(db DatabaseRunner) *Repository {
	cli := gentleman.New()
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)))

	return &Repository{db: db, c: cli, b: ratelimit.NewBucket(10*time.Millisecond, 10)}
}

func (r *Repository) Run(q arangolite.Runnable) (*arangolite.Result, error) {
	r.b.Wait(1)

	result, err := r.db.RunAsync(q)
	if err != nil {
		return nil, stacktrace.Propagate(errs.ErrDatabase, err.Error())
	}

	return result, nil
}

func (r *Repository) RunSync(q arangolite.Runnable) ([]byte, error) {
	r.b.Wait(1)

	result, err := r.db.Run(q)
	if err != nil {
		return nil, stacktrace.Propagate(errs.ErrDatabase, err.Error())
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
		return nil, stacktrace.Propagate(err, "%s %s failed (network)", method, url)
	}

	if !res.Ok {
		errRes := &response{}
		if err := json.Unmarshal(res.Bytes(), errRes); err != nil {
			return nil, stacktrace.Propagate(errs.ErrServiceError, "%s %s returned a %d", method, url, res.StatusCode)
		}

		if strings.Contains(errRes.Raw, "not found") {
			return nil, stacktrace.Propagate(errs.ErrNotFound, "%s %s returned a %d", method, url, res.StatusCode)
		}

		return nil, stacktrace.Propagate(errors.New(errRes.Raw), "%s %s returned a %d", method, url, res.StatusCode)
	}

	return res.Bytes(), nil
}
