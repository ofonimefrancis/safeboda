package promo

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/ofonimefrancis/safeboda/common/log"
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
	Radius         int           `json:"radius"`
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

func (datastore *DatastoreSession) GetEvent(id bson.ObjectId) (Event, error) {
	var event Event
	err := datastore.event().Find(bson.M{"id": id}).One(&event)
	if err != nil {
		return event, err
	}
	return event, nil
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
	err := datastore.promo().Find(bson.M{"code": promo_code, "isactive": true}).One(&promo)
	return err == nil
}

func (datastore *DatastoreSession) IsExpired(code string) bool {
	var promo Promo
	err := datastore.promo().Find(bson.M{"code": code, "isexpired": true}).One(&promo)
	if err != nil {
		return false
	}
	return true
}

func (datastore *DatastoreSession) GetPromo(code string) (Promo, error) {
	var promo Promo
	err := datastore.promo().Find(bson.M{"code": code}).One(&promo)
	if err != nil {
		return promo, err
	}
	return promo, nil
}

func (datastore *DatastoreSession) UpdateRadius(promo Promo, radius int) error {
	return datastore.promo().Update(bson.M{"code": promo.Code}, bson.M{"radius": radius})
}

func (datastore *DatastoreSession) GetAllActivePromos(page string) ([]Promo, error) {
	var promos []Promo
	var pageInt int

	if page == "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Info("Error converting from string to integer. Setting pageInt to default value of 1")
			pageInt = 1
			log.Info(pageInt)
		}
	}

	pageSize := 20
	offset := pageSize * pageInt

	if err := datastore.promo().Find(bson.M{"isactive": true, "isexpired": false}).Skip(offset).Limit(pageSize).All(&promos); err != nil {
		return promos, err
	}
	return promos, nil
}

func (datastore *DatastoreSession) DeactivatePromoCode(code string) error {
	if datastore.PromoAlreadyExists(code) {
		if !datastore.IsActive(code) {
			return errors.New("Code is already deactivated")
		}

		promo, err := datastore.GetPromo(code)
		if err != nil {
			return err
		}

		promo.IsActive = false
		promo.IsExpired = true

		return datastore.promo().Update(bson.M{"code": code}, promo)
	}

	return errors.New("Promo code not found in our system")
}

func (datastore *DatastoreSession) GetAllPromos(page string) ([]Promo, error) {
	var promos []Promo
	var pageInt int

	if page == "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Info("Error converting from string to integer. Setting pageInt to default value of 1")
			pageInt = 1
			log.Info(pageInt)
		}
	}

	pageSize := 20
	offset := pageSize * pageInt //Ideally  pageSize * (pageInt -1)

	err := datastore.promo().Find(bson.M{}).Skip(offset).Limit(pageSize).All(&promos)
	if err != nil {
		return promos, err
	}
	return promos, nil
}
