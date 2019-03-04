package promo

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"googlemaps.github.io/maps"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/ofonimefrancis/safeboda/common"
	polyline "github.com/twpayne/go-polyline"
)

type PromoCodeDetails struct {
	ID             bson.ObjectId `json:"id"`
	Code           string        `json:"code"`
	Radius         int           `json:"radius"`
	Amount         float64       `json:"amount"`
	ExpirationDate time.Time     `json:"expiration_date"`
	IsExpired      bool          `json:"isexpired"`
	IsActive       bool          `json:"isactive"`
	Event          struct {
		ID      bson.ObjectId `json:"event_id"`
		Name    string        `json:"name"`
		Address string        `json:"address"`
		Lat     float64       `json:"lat"`
		Lng     float64       `json:"lng"`
	} `json:"event"`
}

func degreeToRadian(d float64) float64 {
	return d * math.Pi / 180
}

//Distance Returns the distance between two coordinates in meters
func Distance(p1 Coordinate, p2 Coordinate) float64 {
	//Havesine Formula
	var earthRadius float64 = 6378100 // Earth Radius in meters

	var lat1, lat2, long1, long2 float64
	lat1 = degreeToRadian(p1.Lat)
	lat2 = degreeToRadian(p2.Lat)
	long1 = degreeToRadian(p1.Lng)
	long2 = degreeToRadian(p2.Lng)

	diffLat := lat2 - lat1
	diffLon := long2 - long1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRadius
}

func GetAddressCoordinate(address string) Coordinate {
	mapClient := common.GetMapClient()
	geocodingRequest := &maps.GeocodingRequest{Address: address}
	result, err := mapClient.Geocode(context.Background(), geocodingRequest)
	if err != nil {
		log.Fatal(err)
	}

	return Coordinate{Lng: result[0].Geometry.Location.Lng, Lat: result[0].Geometry.Location.Lat}
}

func calculateDistanceToEvent(request validatePromo, originCoord Coordinate, event Event, promo Promo, c *gin.Context) {
	mapClient := common.GetMapClient()

	result, err := mapClient.Geocode(context.Background(), &maps.GeocodingRequest{Address: request.Destination})
	if err != nil {
		log.Fatal(err.Error())
	}

	destinationCoordinate := Coordinate{Lat: result[0].Geometry.Location.Lat, Lng: result[0].Geometry.Location.Lng}
	eventCoordinate := Coordinate{Lat: event.Coordinate.Lat, Lng: event.Coordinate.Lng}

	distance := Distance(eventCoordinate, destinationCoordinate)

	if distance > float64(promo.Radius) {
		log.Println("User is not within the event location")
		c.JSON(http.StatusBadRequest, gin.H{"message": "You are not within the event radius"})
		return
	}

	var coordinates = [][]float64{
		{destinationCoordinate.Lat, destinationCoordinate.Lng},
		{eventCoordinate.Lat, eventCoordinate.Lng},
	}

	polylineBytes := polyline.EncodeCoords(coordinates)
	pr := dumpObject(promo, event)

	c.JSON(http.StatusOK, gin.H{"promo_details": pr, "polyline": string(polylineBytes)})
}

func dumpObject(promo Promo, event Event) PromoCodeDetails {
	var pr PromoCodeDetails
	pr.ID = promo.ID
	pr.Code = promo.Code
	pr.IsActive = promo.IsActive
	pr.ExpirationDate = promo.ExpirationDate
	pr.IsExpired = promo.IsExpired
	pr.Radius = promo.Radius
	pr.Amount = promo.Amount
	pr.Event.ID = event.ID
	pr.Event.Name = event.Name
	pr.Event.Address = event.Address
	pr.Event.Lat = event.Coordinate.Lat
	pr.Event.Lng = event.Coordinate.Lng
	return pr
}
