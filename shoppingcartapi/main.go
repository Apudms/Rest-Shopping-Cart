package main

import (
	"updata/shoppingcartapi/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	// session
	store := session.New()

	app := fiber.New()

	// Middleware to check login
	CheckLogin := func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.SendStatus(200)
		}

		return c.SendStatus(400)
	}

	// controllers
	authAPIController := controllers.InitAuthAPIController(store)
	productAPIController := controllers.InitProductAPIController(store)
	cartAPIController := controllers.InitCartAPIController(store)
	transaksiAPIController := controllers.InitTransaksiAPIController(store)

	app.Post("/login", authAPIController.Login)
	app.Post("/register", authAPIController.Register)
	app.Get("/logout", CheckLogin, authAPIController.Logout)

	prod := app.Group("/api")
	prod.Get("/", productAPIController.GetAllProduct)
	prod.Post("/products/create", CheckLogin, productAPIController.CreateProduct)
	prod.Get("/products/detail/:id", productAPIController.DetailProduct)
	prod.Put("/products/edit/:id", CheckLogin, productAPIController.EditProduct)
	prod.Delete("/products/:id", CheckLogin, productAPIController.DeleteProduct)
	prod.Get("/addtocart/:cartid/product/:productid", CheckLogin, cartAPIController.InsertToCart)

	cart := app.Group("/shoppingcart")
	cart.Get("/:cartid", CheckLogin, cartAPIController.GetShoppingCart)

	transaksi := app.Group("/checkout")
	transaksi.Get("/:userid", CheckLogin, transaksiAPIController.InsertToTransaksi)

	history := app.Group("/history")
	history.Get("/:userid", CheckLogin, transaksiAPIController.GetTransaksi)
	history.Get("/detail/:transaksiid", CheckLogin, transaksiAPIController.DetailTransaksi)

	app.Listen(":3000")
}
