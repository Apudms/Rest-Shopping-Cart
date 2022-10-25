package controllers

import (
	"fmt"
	"strconv"
	"updata/shoppingcartapi/database"
	"updata/shoppingcartapi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type ProductAPIController struct {
	Db    *gorm.DB
	store *session.Store
}

func InitProductAPIController(s *session.Store) *ProductAPIController {
	db := database.InitDb()
	db.AutoMigrate(&models.Product{})

	return &ProductAPIController{Db: db, store: s}
}

// Routing
// GET /products
func (controller *ProductAPIController) GetAllProduct(c *fiber.Ctx) error {
	// Load all Products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.JSON(fiber.Map{
		"status":   200,
		"UserId":   val,
		"Products": products,
	})
}

// POST /products/create
func (controller *ProductAPIController) CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.SendStatus(400)
	}

	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

			// Save the files to disk:
			product.Image = fmt.Sprintf("public/prod-images/%s", file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("public/prod-images/%s", file.Filename)); err != nil {
				return err
			}
		}
	}

	// Save data product
	err := models.CreateProduct(controller.Db, &product)
	if err != nil {
		return c.SendStatus(400)
	}
	// if succeed
	return c.JSON(fiber.Map{
		"status":   200,
		"Products": product,
	})
}

// GET /products/detail:id
func (controller *ProductAPIController) DetailProduct(c *fiber.Ctx) error {
	params := c.AllParams()

	intId, errs := strconv.Atoi(params["id"])

	if errs != nil {
		fmt.Println(errs)
	}

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	return c.JSON(fiber.Map{
		"status":   200,
		"UserId":   val,
		"Products": product,
	})
}

// PUT /products/edit/:id
func (controller *ProductAPIController) EditProduct(c *fiber.Ctx) error {
	var product models.Product

	params := c.AllParams() // "{"id": "1"}"
	intId, _ := strconv.Atoi(params["id"])
	product.Id = intId

	if err := c.BodyParser(&product); err != nil {
		return c.SendStatus(400)
	}

	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			product.Image = fmt.Sprintf("public/prod-images/%s", file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("public/prod-images/%s", file.Filename)); err != nil {
				return err
			}
		}
	}

	// Save product
	err := models.UpdateProduct(controller.Db, &product)
	if err != nil {
		return c.SendStatus(400)
	}

	// if succeed
	return c.JSON(fiber.Map{
		"status":   200,
		"Products": product,
	})
}

// GET /products/hapus/:id
func (controller *ProductAPIController) DeleteProduct(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intId, errs := strconv.Atoi(params["id"])

	if errs != nil {
		fmt.Println(errs)
	}

	var product models.Product
	err := models.DeleteProductById(controller.Db, &product, intId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil dihapus!",
	})
}
