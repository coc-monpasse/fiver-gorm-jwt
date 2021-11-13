package handler

import (
	"fmt"
	"mo_fiber_1/database"
	"mo_fiber_1/model"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))
	if uid != n {
		return false
	}
	return true
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.Find(&user, id)
	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "no user fount with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "user found", "data": user})
}

func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "review your input", "data": err})
	}
	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "coudn't hash the password", "data": err})

	}
	user.Password = hash
	fmt.Println(user.Email)
	db := database.DB

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "couldn't create the user", "data": err})
	}
	// result := db.Create(&user)
	// if result.Error != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "couldn't create the user", "data": err})

	// }

	newuser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}
	return c.JSON(fiber.Map{"status": "success", "message": "user created", "data": newuser})

}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzcwODgxMjAsInVzZXJfaWQiOjIsInVzZXJuYW1lIjoiY29jIn0.c5h_HW9fcWieoHLomFkBkHHj0-RZ9dNWmHOZEjWy8Bs

func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Username string `json:"username"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "review your input", "data": err})
	}
	ids := c.Params("id")
	token := c.Locals("user").(*jwt.Token)
	if !validToken(token, ids) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "review your input", "data": nil})
	}
	id, _ := strconv.ParseFloat(ids, 64)

	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(reflect.TypeOf(claims["user_id"]), reflect.TypeOf(id))
	if id != claims["user_id"] {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "you are trying access other resources", "data": nil})

	}
	db := database.DB
	var user model.User
	db.First(&user, id)
	user.Username = uui.Username
	db.Save(&user)
	return c.JSON(fiber.Map{"status": "success", "message": "user succesfully updated", "data": user})

}
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "review your input", "data": nil})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)
	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "invalid token idt", "data": nil})
	}
	if !validUser(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "not valid user", "data": nil})
	}
	db := database.DB
	var user model.User
	db.First(&user, id)
	db.Delete(&user)
	return c.JSON(fiber.Map{"status": "error", "message": "user succesfully deleted", "data": nil})

}
