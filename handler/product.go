package handler

import (
	"mo_fiber_1/database"
	"mo_fiber_1/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllProducts(c *fiber.Ctx) error {
	db := database.DB
	var products []model.Product
	db.Find(&products)
	return c.JSON(fiber.Map{"status": "success", "message": "all products", "data": products})
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product
	db.Find(&product, id)
	if product.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "success",
			"message": "no product found with that id", "data": nil})
	}
	return c.JSON(fiber.Map{
		"status": "success", "message": "product found", "data": product,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	db := database.DB
	product := new(model.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error", "message": "couldn't create the product", "data": nil,
		})
	}
	db.Create(&product)
	return c.JSON(fiber.Map{
		"status": "success", "message": "product created", "data": product,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var product model.Product
	db.First(&product, id)
	if product.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "success",
			"message": "no product found with that id", "data": nil})
	}
	db.Delete(&product)
	return c.JSON(fiber.Map{
		"status": "success", "message": "product deleted", "data": nil,
	})
}
