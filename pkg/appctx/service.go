package appctx

import (
	"context"

	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type Services struct {
	Database         db.IDBService
	ChangeService    changes.IChangeService
	ReflectorService reflector.IReflectorService
}

type ServiceContext interface {
	Service() *Services
	context.Context
}
