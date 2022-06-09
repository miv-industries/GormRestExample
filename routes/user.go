package routes

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/miv-industries/GormRestExample/database"
	"github.com/miv-industries/GormRestExample/models"
	"github.com/miv-industries/GormRestExample/validators"
)

// Helpers
type User struct {
	// this is not the model User, see this as the serializer
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

// CRUD

func CreateUser(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// field validations because a human will input data here
	if validation_errors := validators.ValidateUser(user); validation_errors != nil {
		return c.Status(400).JSON(validation_errors)
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)
	// in this case User is our serializer
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func GetUser(c *fiber.Ctx) error {
	//this takes an int from path params and puts it in id or returns error
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	// we will filter for this user which is a pointer to an user
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	// we will filter for this user which is a pointer to an user
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateData UpdateUser
	// here we use a pointer to update data because we only want to take the first name and last name of the request
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	fmt.Println(updateData)
	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	// field validations because a human will input data here
	if validation_errors := validators.ValidateUser(user); validation_errors != nil {
		return c.Status(400).JSON(validation_errors)
	}

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	// we will filter for this user which is a pointer to an user and will be filled with the found user data
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// here we attempt to delete the user by passing a pointer to our user
	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully deleted user")

}
