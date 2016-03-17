package controllers

import (
	"net/http"

	"github.com/solher/snakepit"
	"git.wid.la/versatile/versatile-server/models"

	"golang.org/x/net/context"
)

type (
	DashboardsInter interface{}

	DashboardsValidator interface{}

	DashboardsCtrl struct {
		i DashboardsInter
		v DashboardsValidator
		r *snakepit.Render
	}
)

func NewDashboardsCtrl(
	i DashboardsInter,
	v DashboardsValidator,
	r *snakepit.Render,
) *DashboardsCtrl {
	return &DashboardsCtrl{i: i, v: v, r: r}
}

func (c *DashboardsCtrl) Find(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m := models.Dashboard{}
	m.Name = "Name"
	c.r.JSON(w, 200, m)
}
