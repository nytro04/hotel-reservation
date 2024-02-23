package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nytro04/hotel-reservation/api"
	"github.com/nytro04/hotel-reservation/api/middleware"
	"github.com/nytro04/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":9999", "The listen address for the API Server")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// Handler initialization
	var (
		hotelStore   = db.NewMongHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		authApi        = app.Group("/api")
		apiV1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		adminApi       = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// AUTH
	authApi.Post("/auth", authHandler.HandleAuth)

	// USER HANDLERS
	authApi.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// HOTEL HANDLERS
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room handlers
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	// Booking Handlers
	// user handler
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	// admin handler
	adminApi.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
