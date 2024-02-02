package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nytro04/hotel-reservation/api"
	"github.com/nytro04/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":9999", "The listen address for the API Server")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)

	apiV1 := app.Group("/api/v1")

	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	app.Listen(*listenAddr)
}
