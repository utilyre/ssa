package main

import (
	"github.com/utilyre/ssa/internal/database"
	"github.com/utilyre/ssa/internal/handler"
	"github.com/utilyre/ssa/internal/logger"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/store"
	"github.com/utilyre/ssa/internal/view"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			logger.New,
			database.New,
			store.New,
			router.New,
			view.New,
		),
		fx.Invoke(
			handler.HandleStatic,
		),
	).Run()
}
