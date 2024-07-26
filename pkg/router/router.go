package router

import "github.com/gorilla/mux"

// NewRouter builds and returns a new router from routes
func NewRouter() *mux.Router {
	// When StrictSlash == true, if the route path is "/path/", accessing "/path" will perform a redirect to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)
	router.Use(Logger)
	//Path
	sub := router.PathPrefix("/api/v1").Subrouter()
	addRoutes(sub, routes)
	addRoutes(sub, templateRoutes)

	return router
}

func addRoutes(sub *mux.Router, routes Routes) {
	for _, route := range routes {
		sub.
			HandleFunc(route.Pattern, route.HandlerFunc).
			Name(route.Name).
			Methods(route.Method)
	}
}
