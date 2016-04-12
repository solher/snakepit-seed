package app

// func AddRoutes(
// 	r chi.Router,
// 	users chi.HandlerFunc,
// ) {
// 	r.Handle("/users", users)
// }

// func UsersRoutes(router chi.Router) {
//     		router.Route("/", func(router chi.Router) {
// 			router.Post("/", u.Create)
// 			router.Get("/", u.Find)
// 			router.Put("/", u.Update)
// 			router.Delete("/", u.Delete)
// 		})

// 		// CRUD by key operations
// 		router.Route("/:key", func(router chi.Router) {
// 			router.Get("/", u.FindByKey)
// 			router.Put("/", u.UpdateByKey)
// 			router.Delete("/", u.DeleteByKey)
// 		})

// 		// Custom routes
// 		router.Post("/signin", u.Signin)
// 		router.Route("/me", func(router chi.Router) {
// 			router.Get("/", u.FindSelf)
// 			router.Get("/session", u.CurrentSession)
// 			router.Post("/signout", u.Signout)
// 			router.Post("/password", u.UpdateSelfPassword)
// 		})
// }

// func Dispatcher(
// 	vip *viper.Viper,
// 	render *snakepit.Render,
// 	db *arangolite.DB,
// ) func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 		router := chi.NewRouter()

// 		router.Handle("/users", )

// 		router.ServeHTTPC(ctx, w, r)
// 	}
// }
