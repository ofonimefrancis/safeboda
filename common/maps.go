package common

import (
	"log"

	"googlemaps.github.io/maps"
)

func GetMapClient() *maps.Client {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyA8dR77hj6eMWlTpfjZif3pkPmpX0NyIA0"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return c
}
