package route

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"github.com/phassans/frolleague/api"
	"github.com/phassans/frolleague/engines"
)

// APIServerHandler returns a Gzip handler
func APIServerHandler(engines engines.Engine) http.Handler {
	r := newAPIRouter(engines)
	return gziphandler.GzipHandler(r)
}

func newAPIRouter(engines engines.Engine) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", controller.NewRESTRouter(engines))

	return r
}
