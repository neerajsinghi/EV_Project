package router

import 	"net/http"


// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	Route{
		"",
		"GET",
		"/",
		hello,
	},
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
