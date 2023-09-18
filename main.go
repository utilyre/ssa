package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/utilyre/ssa/internal/database"
	"github.com/utilyre/ssa/internal/handler"
	"github.com/utilyre/ssa/internal/logger"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/store"
	"github.com/utilyre/ssa/internal/templates"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			logger.New,
			database.New,
			store.New,
			validator.New,
			router.New,
			templates.New,
		),
		fx.Invoke(
			handler.HandleUsers,
			handler.HandleStatic,
		),
	).Run()
}
