package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/danangkonang/cake-store-restful/config"
	"github.com/danangkonang/cake-store-restful/helper"
	"github.com/danangkonang/cake-store-restful/router"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type loggingRw struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingRw(w http.ResponseWriter) *loggingRw {
	return &loggingRw{w, http.StatusOK}
}

func (w *loggingRw) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		rw := newLoggingRw(w)
		next.ServeHTTP(rw, r)
		rr := fmt.Sprintf("%s %s %s %d %s", r.Host, r.Method, r.RequestURI, rw.statusCode, r.UserAgent())
		helper.LoggerAccees(rr)
	})
}

func Run() {
	r := mux.NewRouter().StrictSlash(false)
	r.Use(loggingHandler)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helper.MakeRespon(w, 404, "page not found", nil)
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helper.MakeRespon(w, http.StatusMethodNotAllowed, "Method NotAllowed", nil)
	})

	router.CakeRouter(r, config.Connection())
	serverloging := fmt.Sprintf("local server started at http://localhost:%s", os.Getenv("APP_PORT"))
	fmt.Println(serverloging)

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), handlers.CORS(
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)(r))
}
