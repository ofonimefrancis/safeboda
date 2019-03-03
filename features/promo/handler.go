package promo

import (
	"context"

	"github.com/ofonimefrancis/safeboda/common/mgo"
)

type Handler struct {
	datastore *Datastore
}

type PromoDataProvider struct {
	datastore *Datastore
}

func NewPromoDataProvider(initContext context.Context, database *mgo.Database) *PromoDataProvider {
	datastore := NewDatastore(initContext, database)
	return &PromoDataProvider{datastore}
}

func NewHandler(initContext context.Context, database *mgo.Database) *Handler {
	datastore := NewDatastore(initContext, database)
	handler := &Handler{datastore}
	return handler
}
