package api

import (
	"context"
	"github.com/hina1314/hina/server/db"
)

type API struct {
	db *db.DB
}

func NewAPI() API {
	ctx := context.Background()
	return API{
		db: db.NewDB(ctx),
	}
}
