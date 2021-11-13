package handler

import (
	"mo_fiber_1/config"
	"mo_fiber_1/database"
	"mo_fiber_1/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).Find(&user).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).Find(&user).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	var ud UserData
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "error", "message": "Error on login request", "data": err},
		)
	}
	identity := input.Identity
	pass := input.Password
	email, err := GetUserByEmail(identity)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"status": "error", "message": "error on email", "data": err},
		)
	}
	// user,err:=GetUserByUsername(identity)
	// if err !=nil{
	// 	return c.Status(fiber.StatusUnauthorized).JSON(
	// 		fiber.Map{"status":"error","message":"error on Username","data":err}
	// 	)
	// }
	if email == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"status": "error", "message": "user not found", "data": err},
		)
	}
	ud = UserData{
		ID:       email.ID,
		Username: email.Username,
		Email:    email.Email,
		Password: email.Password,
	}
	if !CheckPasswordHash(pass, ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"status": "error", "message": "InvalidPassword", "data": nil},
		)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
