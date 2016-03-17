package app

import (
	"github.com/pressly/chi"
	"git.wid.la/versatile/versatile-server/controllers"
)

func AddRoutes(
	r chi.Router,
	uc *controllers.UsersCtrl,
	dc *controllers.DashboardsCtrl,
) {
	// USERS
	r.Route("/users", func(r chi.Router) {
		// CRUD operations
		r.Route("/", func(r chi.Router) {
			r.Post("/", uc.Create)
			r.Get("/", uc.Find)
			r.Put("/", uc.Update)
			r.Delete("/", uc.Delete)
		})

		// CRUD by key operations
		r.Route("/:key", func(r chi.Router) {
			r.Get("/", uc.FindByKey)
			r.Put("/", uc.UpdateByKey)
			r.Delete("/", uc.DeleteByKey)
		})

		// Custom routes
		r.Post("/signin", uc.Signin)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", uc.FindSelf)
			r.Get("/session", uc.CurrentSession)
			r.Post("/signout", uc.Signout)
			r.Post("/password", uc.UpdateSelfPassword)
		})
	})

	r.Route("/dashboards", func(r chi.Router) {
		r.Get("/", dc.Find)
	})
}
