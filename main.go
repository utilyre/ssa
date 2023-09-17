package main

import (
	"github.com/utilyre/ssa/internal/database"
	"github.com/utilyre/ssa/internal/env"
	"github.com/utilyre/ssa/internal/handler"
	"github.com/utilyre/ssa/internal/logger"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/view"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			env.New,
			logger.New,
			database.New,
			router.New,
			view.New,
		),
		fx.Invoke(
			handler.HandleHC,
		),
	).Run()
}
