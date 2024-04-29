package appctx

import (
	"context"

	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
)

type AppContext struct {
	context.Context
	environment *Environment
	service     *Services
}

func Init(
	db db.DBServiceInterface,
	changeService changes.ChangeServiceInterface,
) (AppContext, error) {
	return AppContext{
		Context:     context.Background(),
		environment: &Environment{},
		service: &Services{
			Database:      db,
			ChangeService: changeService,
		},
	}, nil
}

func (self AppContext) Environment() *Environment {
	return self.environment
}

func (self AppContext) Service() *Services {
	return self.service
}
