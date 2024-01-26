package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":9999", "The listen address for the API Server")

	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)

	apiV1.Get("/user", handleUser)

	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working just fine"})
}

func handleUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "user foo"})
}
