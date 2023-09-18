package router

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/utilyre/xmate"
	"go.uber.org/fx"
)

type Router struct {
	router       *mux.Router
	errorHandler xmate.ErrorHandler
}

type MiddlewareFunc func(next http.Handler) xmate.Handler

func (r Router) Subrouter(prefix string) Router {
	return Router{
		router:       r.router.PathPrefix(prefix).Subrouter(),
		errorHandler: r.errorHandler,
	}
}

func (r Router) Use(mwf MiddlewareFunc) {
	r.router.Use(func(next http.Handler) http.Handler {
		return r.errorHandler.Handle(mwf(next))
	})
}

func (r Router) Handle(path string, handler xmate.Handler) *mux.Route {
	return r.router.Handle(path, r.errorHandler.Handle(handler))
}

func (r Router) HandleFunc(path string, handler xmate.HandlerFunc) *mux.Route {
	return r.router.HandleFunc(path, r.errorHandler.HandleFunc(handler))
}

func New(lc fx.Lifecycle, l *slog.Logger) Router {
	r := Router{
		router: mux.NewRouter(),
		errorHandler: func(w http.ResponseWriter, r *http.Request) {
			err := r.Context().Value("error").(error)

			httpErr := new(xmate.HTTPError)
			if !errors.As(err, &httpErr) {
				httpErr.Code = http.StatusInternalServerError
				httpErr.Message = http.StatusText(httpErr.Code)

				l.Warn("HTTP handler failed", "method", r.Method, "path", r.URL.Path, "error", err)
			}

			http.Error(w, httpErr.Message, httpErr.Code)
		},
	}
	srv := http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: r.router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						return
					}

					l.Error("failed to listen and serve", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return r
}
