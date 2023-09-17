package router

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/utilyre/ssa/internal/env"
	"github.com/utilyre/xmate"
	"go.uber.org/fx"
)

type Router struct {
	*mux.Router
	ErrorHandler xmate.ErrorHandler
}

func New(lc fx.Lifecycle, e env.Env, l *slog.Logger) Router {
	r := Router{
		Router: mux.NewRouter(),
		ErrorHandler: func(w http.ResponseWriter, r *http.Request) {
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
		Addr:    ":" + e.BEPort,
		Handler: r,
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
