package controllers

import (
	"updata/shoppingcartapi/database"
	"updata/shoppingcartapi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthAPIController struct {
	Db    *gorm.DB
	store *session.Store
}

func InitAuthAPIController(s *session.Store) *AuthAPIController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &AuthAPIController{Db: db, store: s}
}

// POST /login
func (controller *AuthAPIController) Login(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	var user models.User
	var postlogin LoginForm

	if err := c.BodyParser(&postlogin); err != nil {
		return c.SendStatus(500)
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, postlogin.Username)
	if errs != nil {
		return c.SendStatus(400) // Gagal login (Username tidak ditemukan)
	}

	// Compare password
	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(postlogin.Password))
	if compare == nil { // compare == nil artinya hasil compare di atas true
		sess.Set("username", user.Username)
		sess.Set("userId", user.ID)
		sess.Save()

		return c.JSON(fiber.Map{
			"status":  200,
			"message": "Login berhasil!",
		})
	}

	return c.SendStatus(500)
}

// POST /register
func (controller *AuthAPIController) Register(c *fiber.Ctx) error {
	var user models.User
	var cart models.Cart

	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(400) // Bad Request, // RegistUser is not complete
	}

	// Hash password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	sHash := string(bytes)

	// Simpan hashing, bukan plain passwordnya
	user.Password = sHash

	// Save user
	err := models.CreateUser(controller.Db, &user)
	if err != nil {
		return c.SendStatus(500) // Server error (Gagal menyimpan user)
	}

	// Find user
	errs := models.FindUserByUsername(controller.Db, &user, user.Username)
	if errs != nil {
		return c.SendStatus(500) // Server error (Gagal menyimpan user)
	}

	// Membuat cart berdasarkan data register
	errCart := models.CreateCart(controller.Db, &cart, user.ID)
	if errCart != nil {
		return c.SendStatus(500) // Server error (Gagal menyimpan user)
	}

	// Registrasi berhasil
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil!",
	})
}

// GET /logout
func (controller *AuthAPIController) Logout(c *fiber.Ctx) error {

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Destroy()

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Anda baru saja melakukan Logout!",
	})
}
