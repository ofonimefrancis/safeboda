package db

import (
	"context"

	"github.com/globalsign/mgo/bson"
	"github.com/ofonimefrancis/safeboda/common/log"
	"github.com/ofonimefrancis/safeboda/common/mgo"
)

var collections = []string{
	"promo",
	"event",
}

type Helper struct {
	c       context.Context
	cancel  context.CancelFunc
	session *mgo.Session
	dbName  string
}

func NewHelper(session *mgo.Session, dbName string) *Helper {
	c, cancel := context.WithCancel(context.Background())
	return &Helper{
		session: session,
		// nolint: golint
		c:      context.WithValue(c, mgo.ContextKey, session),
		cancel: cancel,
		dbName: dbName,
	}
}

func (self Helper) Session() *mgo.Session {
	return mgo.FromContext(self.c)
}

func (self Helper) CloseSession() {
	self.cancel()
}

func (self Helper) Cleanup() {
	s := self.Session()

	database := s.DB(self.dbName)
	for _, name := range collections {
		err := removeCollection(database, name)
		if err != nil {
			log.Warningf("db helper cant remove collection %s err=%v", name, err)
		}
	}
}

func (self Helper) RemoveCollection(name string) error {
	s := self.Session()

	database := s.DB(self.dbName)
	return removeCollection(database, name)
}

func removeCollection(database *mgo.Database, name string) error {
	collection := database.C(name)
	_, err := collection.RemoveAll(bson.M{})
	return err
}
