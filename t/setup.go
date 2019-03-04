package t

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ofonimefrancis/safeboda/common/mgo"
	"github.com/ofonimefrancis/safeboda/common/server"
	"github.com/ofonimefrancis/safeboda/features/promo"
)

var (
	counter    int
	s, addr    = server.StartAsync(engine())
	defaultURL = fmt.Sprintf("http://%s", addr)
	mongoHost  = "mongodb://localhost:27017"
	mongoDB    = "safeboda"
)

func engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		counter++
		c.String(http.StatusOK, strconv.Itoa(counter))
	})
	r.GET("redirect", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})
	return r
}

func getDataStoreSession() *promo.DatastoreSession {
	initContext := context.Background()
	datastore := promo.NewDatastore(initContext, mgo.New(mongoHost, mongoDB))
	ds := datastore.OpenSession(initContext)

	return ds
}
