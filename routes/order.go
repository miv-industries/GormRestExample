package routes

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/miv-industries/GormRestExample/database"
	"github.com/miv-industries/GormRestExample/models"
	"github.com/miv-industries/GormRestExample/validators"
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

func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id = ?", id)
	if order.ID == 0 {
		return errors.New("Order does not exist")
	}

	return nil
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
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// field validations because a human will input data here
	if validation_errors := validators.ValidateProduct(product); validation_errors != nil {
		return c.Status(400).JSON(validation_errors)
	}

	// we now create the order in our database
	database.Database.Db.Create(&order)

	// these create functions basically use a struct type as a serializer for the response field or nested object
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}
	database.Database.Db.Find(&orders)
	responseOrders := []Order{}

	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)

		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)

	}
	return c.Status(200).JSON(responseOrders)

}

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var order models.Order

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	database.Database.Db.First(&user, order.UserRefer)
	database.Database.Db.First(&product, order.ProductRefer)
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)

	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}
