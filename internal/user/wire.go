//go:build wireinject
// +build wireinject

package user

import (
	"example.com/m/config"
	"example.com/m/internal/database"
	"github.com/google/wire"
)

func InitializeUserHandler() *Handler {
	wire.Build(
		config.ProvideConfig,
		database.ProvideDatabase,
		ProvideUserService,
		ProvideUserHandler,
	)
	return &Handler{}
}
