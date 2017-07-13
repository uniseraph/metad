package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/zanecloud/metad/opts"
	"net/http"
)

type Handler func(c context.Context, w http.ResponseWriter, r *http.Request)

type Service struct {
	api.AgentServiceRegistration
}

func getServiceJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	client, _ := ctx.Value("consul.client").(*api.Client)

	services, err := client.Agent().Services()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := mux.Vars(r)["name"]
	service, ok := services[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if service.Tags == nil {
		service.Tags = []string{}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)

}
func getServiceList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	client, _ := ctx.Value("consul.client").(*api.Client)

	result, err := client.Agent().Services()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func postServiceCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	req := api.AgentServiceRegistration{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Debugf("postServiceCreate::receving a request :%#v", req)

	client, ok := ctx.Value("consul.client").(*api.Client)

	if !ok {
		logrus.Errorf("postServiceCreate:can't get a consul client,ctx is %#v", ctx)
		http.Error(w, fmt.Sprintf("postServiceCreate:can't get a consul client,ctx is %#v", ctx), http.StatusBadRequest)
		return
	}

	if err := client.Agent().ServiceRegister(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

var routers = map[string]map[string]Handler{
	"HEAD": {},
	"GET": {
		"/services/list":              getServiceList,
		"/services/{name:.*}/inspect": getServiceJSON,
	},
	"POST": {
		"/services/create": postServiceCreate,
	},
	"PUT":     {},
	"DELETE":  {},
	"OPTIONS": {},
}

func RunMetad(opts opts.MetadOptions) {

	consulClient, err := api.NewClient(&api.Config{Address: opts.Consul})
	if err != nil {
		logrus.Fatalf("create a consul client error:%s", err.Error())
		return
	}

	r := mux.NewRouter()
	for method, mappings := range routers {
		for route, fct := range mappings {
			logrus.WithFields(logrus.Fields{"method": method, "route": route}).Debug("Registering HTTP route")

			localRoute := route
			localFct := fct
			wrap := func(w http.ResponseWriter, req *http.Request) {
				logrus.WithFields(logrus.Fields{"method": req.Method, "uri": req.RequestURI}).Debug("HTTP request received")

				ctx := context.WithValue(req.Context(), "metad.opts", opts)
				ctx = context.WithValue(ctx, "consul.client", consulClient)

				localFct(ctx, w, req)
			}
			localMethod := method

			//r.Path("/v{version:[0-9.]+}" + localRoute).Methods(localMethod).HandlerFunc(wrap)
			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
		}
	}

	srv := http.Server{
		Handler: r,
		Addr:    opts.Address,
	}

	if err := srv.ListenAndServe(); err != nil {
		logrus.Errorf("run metad err:%s", err.Error())
	}
}
