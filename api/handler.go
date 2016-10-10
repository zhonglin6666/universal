package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
)

const (
	RVERSION = "/v1/"
)

var MasterRouter = map[string]map[string]func(http.ResponseWriter, *http.Request){
	"GET": {
		"db/master/test": getMasterTest,
	},

	"POST": {},

	"DELETE": {},
}

func HandlerRouter(mRouter map[string]map[string]func(http.ResponseWriter, *http.Request)) http.Handler {
	router := mux.NewRouter()

	for method, routes := range mRouter {
		for route, fct := range routes {
			localRoute := RVERSION + route
			localMethod := method
			handler := handlerRequest(fct)
			router.Path(localRoute).Methods(localMethod).HandlerFunc(handler)
		}
	}
	router.NotFoundHandler = http.NotFoundHandler()

	return router
}

func handlerRequest(handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("server API handRequest: %s", r.URL)
		handlerFunc(w, r)
	}
}
