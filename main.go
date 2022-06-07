package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/miv-industries/GormRestExample/database"
	"github.com/miv-industries/GormRestExample/routes"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my awesome API")
}

func setuproutes(app *fiber.App) {
	// welcome endpoint
	app.Get("/api", welcome)

	// User endpoints
	app.Post("api/users", routes.CreateUser)

}

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Get("/api", welcome)

	log.Fatal(app.Listen(":3000"))

}
