package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nytro04/hotel-reservation/api"
	"github.com/nytro04/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoDBEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBEndpoint))
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
		apiV1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		adminApi       = apiV1.Group("/admin", api.AdminAuth)
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
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room handlers
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	// Booking Handlers
	// user handler
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	// this works, but i will prefer to send a PATCH request with cancel params set to true instead of a get
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handler
	adminApi.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")

	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
