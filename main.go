package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/cli"
	"github.com/ofonimefrancis/safeboda/common"
	"github.com/ofonimefrancis/safeboda/common/config"
	"github.com/ofonimefrancis/safeboda/common/log"
	"github.com/ofonimefrancis/safeboda/common/mgo"
	"github.com/ofonimefrancis/safeboda/features/promo"
)

func main() {
	cli.Run(new(config.PackageFlag), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*config.PackageFlag)
		initContext, finishInit := context.WithCancel(context.Background())

		r := gin.Default()
		gin.SetMode(gin.ReleaseMode)

		r.Use(common.EnsureHTTPVersion())
		r.Use(common.SecureHeaders())
		r.Use(common.SilenceSomePanics())

		database := mgo.New(argv.DBHost, argv.DBName)
		r.Use(mgo.DBConnectionMiddleware(database))

		log.Info("Registering features...")
		promoHandler := promo.NewHandler(initContext, database)
		promoFacade := promo.NewFacade(promoHandler)
		promoFacade.RegisterRoute(r.Group(promo.BasePath))
		finishInit()

		return http.ListenAndServe(fmt.Sprintf(":%d", argv.Port), r)
	})
}
