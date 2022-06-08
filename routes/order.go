package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/miv-industries/GormRestExample/database"
	"github.com/miv-industries/GormRestExample/models"
)

/*Example Order {
	id: 1,
	user: {
		id: 23,
		first_name: "Bob",
		last_name: "Marley"
	},
	product: {
		id: 24,
		name: "Macbook",
		serial_number: "3212312"
	}
} */

type Order struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"order_date"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

func CreateOrder(c *fiber.Ctx) error {

	// we check for field errors here and we parse the order into the var
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// we check that the user foregin key object actually exists and place it in the var
	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// we check that the product foreign key object actually exists and place it in the var
	var product models.Product
	if err := findProduct(order.UserRefer, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// we now create the order in our database
	database.Database.Db.Create(&order)

	// these create functions basically use a struct type as a serializer for the response field or nested object
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}
