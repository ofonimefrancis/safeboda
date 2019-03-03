package promo

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/ofonimefrancis/safeboda/common/mgo"
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Event struct {
	ID         bson.ObjectId `json:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	Coordinate Coordinate    `json:"coordinate"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type Promo struct {
	ID             bson.ObjectId `json:"_id"`
	Code           string        `json:"code"`
	Radius         string        `json:"radius"`
	Amount         float64       `json:"amount"`
	ExpirationDate time.Time     `json:"expiration_date"`
	IsExpired      bool          `json:"is_expired"`
	EventID        bson.ObjectId `json:"event_id"`
	IsActive       bool          `json:"is_active"`
}

type Datastore struct {
	database *mgo.Database
}

type DatastoreSession struct {
	database *mgo.Database
}

func NewDatastore(initContext context.Context, database *mgo.Database) *Datastore {
	datastore := &Datastore{
		database: database,
	}
	session := datastore.OpenSession(initContext)
	mgo.EnsureOrUpgradeIndexKey(session.promo(), "code")

	return datastore
}

func (ds *Datastore) OpenSession(c context.Context) *DatastoreSession {
	db := ds.database
	return &DatastoreSession{
		database: db.FromContext(c),
	}
}

func (datastore *DatastoreSession) promo() *mgo.Collection {
	return datastore.database.C("promo")
}

func (datastore *DatastoreSession) event() *mgo.Collection {
	return datastore.database.C("event")
}

func (datastore *DatastoreSession) EventAlreadyExists(name string) bool {
	var event Event
	err := datastore.event().Find(bson.M{"name": name}).One(&event)
	if err != nil {
		return false
	}
	return true
}

func (datastore *DatastoreSession) PromoAlreadyExists(code string) bool {
	var promo Promo
	err := datastore.promo().Find(bson.M{"code": code}).One(&promo)
	if err != nil {
		return false
	}
	return true
}

func (datastore *DatastoreSession) PromoLinkedToEvent(eventObjectID bson.ObjectId) bool {
	var promo Promo
	if err := datastore.promo().Find(bson.M{"eventid": eventObjectID}).One(&promo); err != nil {
		return false
	}
	return true
}

func (datastore *DatastoreSession) NewEvent(event Event) error {
	return datastore.event().Insert(event)
}

func (datastore *DatastoreSession) NewPromo(promo Promo) error {
	return datastore.promo().Insert(promo)
}

func (datastore *DatastoreSession) IsActive(promo_code string) bool {
	var promo Promo
	err := datastore.promo().Find(bson.M{"code": promo_code, "is_active": true}).One(&promo)
	return err == nil
}

func (datastore *DatastoreSession) IsExpired(code string) bool {
	var promo Promo
	err := datastore.promo().Find(bson.M{"code": code, "is_expired": true}).One(&promo)
	if err != nil {
		return false
	}
	return true
}
