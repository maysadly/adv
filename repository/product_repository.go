package repository

import "FoodStore-AdvProg2/domain"

type ProductRepository interface {
	Save(product domain.Product) error
	FindByID(id string) (domain.Product, error)
	Update(id string, product domain.Product) error
	Delete(id string) error
    FindAll() ([]domain.Product, error)
}
