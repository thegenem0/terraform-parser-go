package appctx

import (
	"context"

	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/db"
)

type Services struct {
	Database      db.DBServiceInterface
	ChangeService changes.ChangeServiceInterface
}

type ServiceContext interface {
	Service() *Services
	context.Context
}
