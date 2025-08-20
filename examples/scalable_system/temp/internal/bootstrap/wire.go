//go:build wireinject
// +build wireinject

package bootstrap

import (
	ctrlUser "github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/controller/rest/user"
	repoUser "github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/repository/dummy/user"
	svcUser "github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/service/user"
	"github.com/google/wire"
	redisv9 "github.com/redis/go-redis/v9"

	"github.com/ciazhar/go-zhar/pkg/validator"
)

// Provider set for the user stack
var userSet = wire.NewSet(
	repoUser.NewUserRepository,
	svcUser.NewUserService,
	ctrlUser.NewUserController,
)

// Build the REST module given infra deps from main (validator + redis client).
func InitializeRESTModule(v validator.Validator, rdb *redisv9.Client) *RESTModule {
	wire.Build(
		userSet,
		NewRESTModule,
	)
	return nil
}
