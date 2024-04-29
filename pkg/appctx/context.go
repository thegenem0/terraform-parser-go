package appctx

import (
	"context"

	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type AppContext struct {
	context.Context
	environment *Environment
	service     *Services
}

func Init(
	db db.IDBService,
	changeService changes.IChangeService,
	reflectorService reflector.IReflectorService,
) (AppContext, error) {
	return AppContext{
		Context:     context.Background(),
		environment: &Environment{},
		service: &Services{
			Database:         db,
			ChangeService:    changeService,
			ReflectorService: reflectorService,
		},
	}, nil
}

func (self AppContext) Environment() *Environment {
	return self.environment
}

func (self AppContext) Service() *Services {
	return self.service
}
