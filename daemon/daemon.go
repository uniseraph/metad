package daemon

import (
	"github.com/gorilla/mux"
	"github.com/zanecloud/metad/opts"
	"github.com/Sirupsen/logrus"
	"net/http"
	"context"
)

type Handler func(c context.Context, w http.ResponseWriter, r *http.Request)


var routers = map[string]map[string]Handler{
	"HEAD": {},
	"GET": {
	},
	"POST": {
	},
	"PUT": {},
	"DELETE": {
	},
	"OPTIONS": {
	},
}

func RunMetad(opts opts.MetadOptions, address string) {


	r := mux.NewRouter()
	for method, mappings := range routers {
		for route, fct := range mappings {
			logrus.WithFields(logrus.Fields{"method": method, "route": route}).Debug("Registering HTTP route")

			localRoute := route
			localFct := fct
			wrap := func(w http.ResponseWriter, r *http.Request) {
				logrus.WithFields(logrus.Fields{"method": r.Method, "uri": r.RequestURI}).Debug("HTTP request received")

				ctx := context.WithValue(r.Context(), "metad.opts",opts)

				localFct(ctx, w, r)
			}
			localMethod := method

			r.Path("/v{version:[0-9.]+}" + localRoute).Methods(localMethod).HandlerFunc(wrap)
			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
		}
	}


}
