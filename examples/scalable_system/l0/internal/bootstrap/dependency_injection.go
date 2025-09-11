//go:build wireinject
// +build wireinject

package bootstrap

import (
	ctrlUser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/controller/rest/user"
	repoUser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user"
	svcUser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/service/user"
	"github.com/ciazhar/go-zhar/pkg/validator"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitializeRESTModule Build the REST module given infra deps from main (validator + redis client).
func InitializeRESTModule(v validator.Validator, pool *pgxpool.Pool) *RESTModule {
	wire.Build(
		//repository
		repoUser.NewUserRepository,

		//service
		svcUser.NewUserService,

		//controller
		ctrlUser.NewUserController,

		//rest module
		NewRESTModule,
	)
	return nil
}
