package controllers

import (
	"strconv"
	"updata/shoppingcartapi/database"
	"updata/shoppingcartapi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type TransaksiAPIController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitTransaksiAPIController(s *session.Store) *TransaksiAPIController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Transaksi{})

	return &TransaksiAPIController{Db: db, store: s}
}

// GET /checkout/:userid
func (controller *TransaksiAPIController) InsertToTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksi models.Transaksi
	var cart models.Cart

	err := models.ReadAllProductsInCart(controller.Db, &cart, intUserId)
	if err != nil {
		return c.SendStatus(500)
	}

	errs := models.CreateTransaksi(controller.Db, &transaksi, uint(intUserId), cart.Products)
	if errs != nil {
		return c.SendStatus(500)
	}

	// Delete products in cart
	errss := models.UpdateCart(controller.Db, cart.Products, &cart, uint(intUserId))
	if errss != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Checkout berhasil!",
	})
}

// GET /historytransaksi/:userid
func (controller *TransaksiAPIController) GetTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksis []models.Transaksi
	err := models.ReadTransaksiById(controller.Db, &transaksis, intUserId)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"status":     200,
		"Transaksis": transaksis,
	})

}

// GET /history/detail/:transaksiid
func (controller *TransaksiAPIController) DetailTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intTransaksiId, _ := strconv.Atoi(params["transaksiid"])

	var transaksi models.Transaksi
	err := models.ReadAllProductsInTransaksi(controller.Db, &transaksi, intTransaksiId)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"status":   200,
		"Products": transaksi.Products,
	})
}
