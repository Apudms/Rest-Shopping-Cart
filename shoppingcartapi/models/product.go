package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id         int          `json: "id" validate:"required"`
	Name       string       `json: "name" validate:"required"`
	Image      string       `json: "image" validate:"required"`
	Quantity   int          `json: "quantity" validate:"required"`
	Price      float32      `json: "price" validate:"required"`
	Carts      []*Cart      `gorm:"many2many:cart_products;"`
	Transaksis []*Transaksi `gorm:"many2many:transaksi_products;"`
}

// CREATE
func CreateProduct(db *gorm.DB, newProduct *Product) (err error) {
	err = db.Create(newProduct).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL
func ReadProducts(db *gorm.DB, products *[]Product) (err error) {
	err = db.Find(products).Error
	if err != nil {
		return err
	}
	return nil
}

// READ By Id
func ReadProductById(db *gorm.DB, product *Product, id int) (err error) {
	err = db.Where("id=?", id).First(product).Error
	if err != nil {
		return err
	}
	return nil
}

// UPDATE
func UpdateProduct(db *gorm.DB, product *Product) (err error) {
	db.Save(product)

	return nil
}

// DELETE
func DeleteProductById(db *gorm.DB, product *Product, id int) (err error) {
	db.Where("id=?", id).Delete(product)

	return nil
}
