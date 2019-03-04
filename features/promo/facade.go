package promo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
	ID        bson.ObjectId `json:"_id"`
	Name      string        `json:"name"`
	Address   string        `json:"address"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	CreatedAt time.Time     `json:"created_at"`
}

type promoPayload struct {
	ID     bson.ObjectId `json:"_id"`
	Code   string        `json:"code"`
	Radius float64       `json:"radius"`
	Amount float64       `json:"amount"`
}

type validatePromo struct {
	Code        string `json:"code"`
	Destination string `json:"destination"`
	Origin      string `json:"origin"`
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
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Service Alive")
	})

	r.GET("/all", func(c *gin.Context) {
		page := c.Query("page")

		ds := facade.promoHandler.datastore.OpenSession(context.Background())

		promos, err := ds.GetAllPromos(page)
		if err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error retrieving promos"})
			return
		}
		if len(promos) == 0 {
			promos = []Promo{}
		}
		c.JSON(http.StatusOK, promos)
	})

	r.GET("/active", func(c *gin.Context) {
		ds := facade.promoHandler.datastore.OpenSession(context.Background())
		page := c.Query("page")

		activePromos, err := ds.GetAllActivePromos(page)
		if err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error retrieving promos"})
			return
		}
		if len(activePromos) == 0 {
			activePromos = []Promo{}
		}
		c.JSON(http.StatusOK, activePromos)
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
		event.Coordinate.Lat = requestBody.Latitude
		event.Coordinate.Lng = requestBody.Longitude
		event.CreatedAt = time.Now()

		if err := ds.NewEvent(event); err != nil {
			log.Info("Error adding a new event")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error creating a new event"})
			return
		}
		c.JSON(http.StatusOK, eventResponse{Event: event})
		return
	})

	//Adds a new Promo for an event with (event_id)
	r.POST("/new/:event_id", func(c *gin.Context) {
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
			promo.Radius = requestBody["radius"].(int)
			promo.Amount = requestBody["amount"].(float64)
			promo.ExpirationDate, _ = common.ParseTime(requestBody["expiration_date"].(string)) //Treat err
			promo.ID = bson.NewObjectId()
			promo.IsActive = true
			promo.IsExpired = false
			promo.EventID = eventObjectID
			promo.Code = fmt.Sprintf("SAFE-%s", common.GenerateRandomToken())
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

	r.GET("/deactivate/:code", func(c *gin.Context) {
		code := c.Param("code")
		ds := facade.promoHandler.datastore.OpenSession(context.Background())
		if err := ds.DeactivatePromoCode(code); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Promo code deactivated"})
	})

	//Validate the code, origin and destination
	r.POST("/validate", func(c *gin.Context) {
		ds := facade.promoHandler.datastore.OpenSession(context.Background())

		var requestBody validatePromo

		if err := c.Bind(&requestBody); err != nil {
			log.Info("[Validate Promo Code] Error decoding json into payload struct")
			c.JSON(http.StatusBadRequest, common.ErrSomethingWentWrong)
			return
		}

		if !ds.PromoAlreadyExists(requestBody.Code) {
			log.Info("Promo Code doesn't exist..")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Promo code doesn't exist."})
			return
		}

		if !ds.IsActive(requestBody.Code) {
			log.Info("Promo code isn't active.")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Promo code is inactive"})
			return
		}

		promo, err := ds.GetPromo(requestBody.Code)
		if err != nil {
			log.Info("Can't find promo with specified code")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Can't find Promo with that code"})
			return
		}

		event, err := ds.GetEvent(promo.EventID)
		if err != nil {
			log.Info("Error retrieving event")
			log.Info(event)
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error retrieving event"})
			return
		}

		originCoordinates := GetAddressCoordinate(requestBody.Origin)

		calculateDistanceToEvent(requestBody, originCoordinates, event, promo, c)
	})

	//Configure the radius of an event
	r.POST("/configure/radius", func(c *gin.Context) {
		var requestBody map[string]string
		if err := c.Bind(&requestBody); err != nil {
			log.Info("Error decoding json into struct")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Something went wrong"})
			return
		}

		newRadius, err := strconv.Atoi(requestBody["radius"])
		if err != nil {
			log.Info("Error converting from string to integer")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Something went wrong"})
			return
		}
		promoCode := requestBody["code"]
		ds := facade.promoHandler.datastore.OpenSession(context.Background())
		promo, err := ds.GetPromo(promoCode)
		if err != nil {
			log.Info("Error retrieving promo with the specified code")
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Error retrieving promo"})
			return
		}

		if !promo.IsActive || promo.IsExpired {
			log.Info("Invalid Promo Code")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Promo code is either Inactive or Expired"})
			return
		}

		err = ds.UpdateRadius(promo, newRadius)
		if err != nil {
			log.Info("Error updating radius")
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "Error occurred while updating radius"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Radius updated"})
	})
}
