//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	redisv9 "github.com/redis/go-redis/v9"

	ctrlUser "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/controller/rest/user"
	repoUser "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/repository/dummy/user"
	svcUser "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/service/user"
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
