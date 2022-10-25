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
	// CheckLogin := func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.SendStatus(200)
	// 	}

	// 	return c.SendStatus(400)
	// }

	// controllers
	authAPIController := controllers.InitAuthAPIController(store)
	productAPIController := controllers.InitProductAPIController(store)
	cartAPIController := controllers.InitCartAPIController(store)
	transaksiAPIController := controllers.InitTransaksiAPIController(store)

	app.Post("/login", authAPIController.Login)
	app.Post("/register", authAPIController.Register)
	app.Get("/logout", authAPIController.Logout)

	prod := app.Group("/api")
	prod.Get("/", productAPIController.GetAllProduct)
	prod.Post("/products/create", productAPIController.CreateProduct)
	prod.Get("/products/detail/:id", productAPIController.DetailProduct)
	prod.Put("/products/edit/:id", productAPIController.EditProduct)
	prod.Delete("/products/:id", productAPIController.DeleteProduct)
	prod.Get("/addtocart/:cartid/product/:productid", cartAPIController.InsertToCart)

	cart := app.Group("/shoppingcart")
	cart.Get("/:cartid", cartAPIController.GetShoppingCart)

	transaksi := app.Group("/checkout")
	transaksi.Get("/:userid", transaksiAPIController.InsertToTransaksi)

	history := app.Group("/history")
	history.Get("/:userid", transaksiAPIController.GetTransaksi)
	history.Get("/detail/:transaksiid", transaksiAPIController.DetailTransaksi)

	app.Listen(":3000")
}
