package config

import (
	"github.com/mkideal/cli"
)

type PackageFlag struct {
	cli.Helper
	Port   int    `cli:"p, port" usage:"Application is running on this port" dft:"5000"`
	DBHost string `cli:"db-host" usage:"mongoDB host" dft:"mongodb://localhost:27017"`
	DBName string `cli:"db-name" usage:"mongoDB name" dft:"safeboda"`
}
