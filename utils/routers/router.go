package routers

import (
	"ws-bigiot-services/services"
	"ws-bigiot-services/utils"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type RouteWS struct {
	Name                      string
	Method                    string
	Pattern                   string
	AuthenticationHandlerFunc negroni.HandlerFunc
	HandlerFunc               negroni.HandlerFunc
}

type RoutesWS []RouteWS

var ArrayRoutesWS = RoutesWS{
	/*
			//Example
		RouteWS{
			"Request Name: Login",
			"Request Type: GET/POST/...",
			"WS URL: /webservice/login",
			"Authentiction (Not Required): nil/authentication.RequireTokenAuthentication",
			"Service Functionality: services.Login",
		},
	*/
}

var RoutesWSArray = RoutesWS{
	RouteWS{
		"POST GetRoutes",
		"POST",
		"/ws-bigiot-services/routes",
		nil,
		services.GetRoutes,
	},
	/*
		RouteWS{
			"GET GetRoutes",
			"GET",
			"/ws-bigiot-services/routes",
			nil,
			services.GetRoutes,
		},
	*/
}

var routesIsAliveWSArray = RoutesWS{

	/*
		isAlive?
		method: GET
	*/
	RouteWS{
		"isalive",
		"GET",
		"/webservice/isalive",
		nil,
		utils.IsAlive,
	},

	/*
		isAlive?
		method: POST
	*/
	RouteWS{
		"isalive",
		"POST",
		"/webservice/isalive",
		nil,
		utils.IsAlive,
	},
}

var routesNoPrefixWSArray = RoutesWS{
	RouteWS{
		"RootCheck",
		"GET",
		"/",
		nil,
		utils.IsAlive,
	},
	RouteWS{
		"RootCheck",
		"POST",
		"/",
		nil,
		utils.IsAlive,
	},
}

var includeMap = make(map[string]string)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	addRoutes(router, ArrayRoutesWS)

	return router
}

// ----------------------------------------------------------------------
// addRoutes
// ----------------------------------------------------------------------
func addRoutes(router *mux.Router, AddArray RoutesWS) {
	for _, route := range AddArray {
		var handler negroni.Handler

		handler = route.HandlerFunc
		handler = utils.Logger(handler, route.Name)

		var handlers *negroni.Negroni
		if route.AuthenticationHandlerFunc != nil {
			handlers = negroni.New(
				negroni.HandlerFunc(route.AuthenticationHandlerFunc),
				handler,
			)
		} else {
			handlers = negroni.New(handler)
		}

		router.Handle(route.Pattern, handlers).Methods(route.Method)
	}
}
