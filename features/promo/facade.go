package promo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/ofonimefrancis/safeboda/common"
	"github.com/ofonimefrancis/safeboda/common/log"
)

const (
	BasePath = "promo"
)

func BuildPath(path string) string {
	return fmt.Sprintf("/%s/%s", BasePath, path)
}

type facade struct {
	promoHandler *Handler
}

type eventPayload struct {
	ID         bson.ObjectId `json:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	Coordinate Coordinate    `json:"coordinate"`
	CreatedAt  time.Time     `json:"created_at"`
}

type promoPayload struct {
	ID     bson.ObjectId `json:"_id"`
	Code   string        `json:"code"`
	Radius float64       `json:"radius"`
	Amount float64       `json:"amount"`
}

type eventResponse struct {
	Event Event `json:"event"`
}

type promoResponse struct {
	Promo Promo `json:"promo"`
}

func NewFacade(handler *Handler) *facade {
	return &facade{handler}
}

func (facade *facade) RegisterRoute(r *gin.RouterGroup) {
	r.GET(BasePath, func(c *gin.Context) {
		c.String(http.StatusOK, "Service Alive")
	})

	r.GET("/promo/deactivate/:code", func(c *gin.Context) {
		//Check if the code already exist and tied to an event
	})

	r.POST("/event/new", func(c *gin.Context) {
		var requestBody eventPayload
		if err := c.Bind(&requestBody); err != nil {
			log.Info("[New Event] Error decoding json into payload struct")
			c.JSON(http.StatusBadRequest, common.ErrSomethingWentWrong)
			return
		}

		ds := facade.promoHandler.datastore.OpenSession(context.Background())

		if ds.EventAlreadyExists(requestBody.Name) {
			log.Info("Event already exists..")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Event Already exists."})
			return
		}

		//Create a new Event
		var event Event
		event.ID = bson.NewObjectId()
		event.Name = requestBody.Name
		event.Address = requestBody.Address
		event.Coordinate = requestBody.Coordinate
		event.CreatedAt = time.Now()

		if err := ds.NewEvent(event); err != nil {
			log.Info("Error adding a new event")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error creating a new event"})
			return
		}
		c.JSON(http.StatusOK, eventResponse{Event: event})
		return
	})

	r.POST("/promo/new/:event_id", func(c *gin.Context) {
		//Add a new Promo
		eventParam := c.Param("event_id")
		var eventObjectID bson.ObjectId
		var requestBody map[string]interface{}

		if err := c.Bind(&requestBody); err != nil {
			log.Info("[New PromoCode] Error decoding json into payload struct")
			c.JSON(http.StatusBadRequest, common.ErrSomethingWentWrong)
			return
		}

		ds := facade.promoHandler.datastore.OpenSession(context.Background())
		if !bson.IsObjectIdHex(eventParam) {
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid Event ID"})
			return
		}
		eventObjectID = bson.ObjectIdHex(eventParam)

		var promo Promo
		if !ds.PromoLinkedToEvent(eventObjectID) {
			promo.Radius = requestBody["radius"].(string)
			promo.Amount = requestBody["amount"].(float64)
			promo.ExpirationDate, _ = common.ParseTime(requestBody["expiration_date"].(string)) //Treat err
			promo.ID = bson.NewObjectId()
			promo.IsActive = true
			promo.IsExpired = false
			promo.EventID = eventObjectID
			promo.Code = common.GenerateRandomToken()
			if err := ds.NewPromo(promo); err != nil {
				log.Info("Error creating a new promo code")
				log.Info(err)
				c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error creating new promo code"})
				return
			}
			c.JSON(http.StatusOK, promoResponse{Promo: promo})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "There exists a promo code for this event"})
		return
	})
}
